
from typing import Iterator
import math
import random
import numpy as np




def generate_blue_noise_signal(samples: int) -> np.ndarray:
    """Generates a 1D audio-style Blue Noise signal using spectral shaping."""
    # White noise baseline input in time domain
    white_noise = np.random.normal(0, 1, samples)
    
    # Transition directly to the frequency spectrum
    spectrum = np.fft.rfft(white_noise)
    frequencies = np.fft.rfftfreq(samples)
    
    # Scale frequency amplitudes by sqrt(f) to achieve +3dB/octave slope
    # Add a tiny offset to avoid dividing or multiplying by zero at index 0
    shaper = np.sqrt(frequencies)
    scaled_spectrum = spectrum * shaper
    
    # Return cleanly back to the structural time-domain signal
    blue_noise = np.fft.irfft(scaled_spectrum, n=samples)
    
    # Normalize to unit standard deviation
    return blue_noise / np.std(blue_noise)

def generate_blue_points(width: float, height: float, r: float, k: int = 30) -> list[tuple[float, float]]:
    """Generates 2D Blue Noise coordinates using Bridson's Poisson-Disk algorithm."""
    # Step 0: Calculate background spatial grid cell boundaries
    cell_size = r / math.sqrt(2)
    grid_w = math.ceil(width / cell_size)
    grid_h = math.ceil(height / cell_size)
    
    # Grid initialized to None (empty space tracking structure)
    grid = [[None for _ in range(grid_h)] for _ in range(grid_w)]
    
    points = []
    active_list = []
    
    # Step 1: Select an initial seed point randomly from the domain space
    p0 = (random.uniform(0, width), random.uniform(0, height))
    points.append(p0)
    active_list.append(p0)
    
    gx = int(p0[0] / cell_size)
    gy = int(p0[1] / cell_size)
    grid[gx][gy] = p0
    
    # Step 2: Iterate through active sample candidates
    while active_list:
        idx = random.randint(0, len(active_list) - 1)
        source = active_list[idx]
        found_candidate = False
        
        # Try generating up to k candidates around our active point source
        for _ in range(k):
            # Sample point inside an annular ring between radius r and 2r
            angle = random.uniform(0, 2 * math.pi)
            distance = random.uniform(r, 2 * r)
            
            cx = source[0] + distance * math.cos(angle)
            cy = source[1] + distance * math.sin(angle)
            
            # Check bounds boundaries
            if not (0 <= cx < width and 0 <= cy < height):
                continue
                
            cell_x = int(cx / cell_size)
            cell_y = int(cy / cell_size)
            
            # Search spatial neighborhood grid for conflicts
            too_close = False
            start_x = max(0, cell_x - 2)
            end_x = min(grid_w - 1, cell_x + 2)
            start_y = max(0, cell_y - 2)
            end_y = min(grid_h - 1, cell_y + 2)
            
            for x in range(start_x, end_x + 1):
                for y in range(start_y, end_y + 1):
                    neighbor = grid[x][y]
                    if neighbor:
                        dist = math.hypot(cx - neighbor[0], cy - neighbor[1])
                        if dist < r:
                            too_close = True
                            break
                if too_close:
                    break
            
            # Step 3: If candidate qualifies, commit it to tracking arrays
            if not too_close:
                new_point = (cx, cy)
                points.append(new_point)
                active_list.append(new_point)
                grid[cell_x][cell_y] = new_point
                found_candidate = True
                break
                
        if not found_candidate:
            active_list.pop(idx)
            
    return points

# Example Usage:
# Find blue points within a 100x100 space, keeping them at least 4 units apart
# blue_coordinates = generate_blue_points(100.0, 100.0, r=4.0)
def stream_blue_points(width: float, height: float, r: float, k: int = 30) -> Iterator[tuple[float, float]]:
    """Yields 2D Blue Noise coordinates one-by-one using an active stream queue."""
    cell_size = r / math.sqrt(2)
    grid_w = math.ceil(width / cell_size)
    grid_h = math.ceil(height / cell_size)
    
    grid = [[None for _ in range(grid_h)] for _ in range(grid_w)]
    active_list = []
    
    # Generate and yield initial seed point
    p0 = (random.uniform(0, width), random.uniform(0, height))
    grid[int(p0[0] / cell_size)][int(p0[1] / cell_size)] = p0
    active_list.append(p0)
    yield p0
    
    while active_list:
        idx = random.randint(0, len(active_list) - 1)
        source = active_list[idx]
        found_candidate = False
        
        for _ in range(k):
            angle = random.uniform(0, 2 * math.pi)
            distance = random.uniform(r, 2 * r)
            cx = source[0] + distance * math.cos(angle)
            cy = source[1] + distance * math.sin(angle)
            
            if not (0 <= cx < width and 0 <= cy < height):
                continue
                
            cell_x = int(cx / cell_size)
            cell_y = int(cy / cell_size)
            
            too_close = False
            for x in range(max(0, cell_x - 2), min(grid_w - 1, cell_x + 2) + 1):
                for y in range(max(0, cell_y - 2), min(grid_h - 1, cell_y + 2) + 1):
                    neighbor = grid[x][y]
                    if neighbor and math.hypot(cx - neighbor[0], cy - neighbor[1]) < r:
                        too_close = True
                        break
                if too_close:
                    break
            
            if not too_close:
                new_point = (cx, cy)
                grid[cell_x][cell_y] = new_point
                active_list.append(new_point)
                found_candidate = True
                yield new_point
                break
                
        if not found_candidate:
            active_list.pop(idx)


def stream_blue_signal(beta: float = 0.5) -> Iterator[float]:
    """Generates an infinite stream of 1D scalar Blue Noise values.
    
    Uses a single-pole high-pass filtering strategy on a white-noise source.
    Beta controls the high-pass cutoff (0.0 to 1.0). Higher balances closer to blue.
    """
    last_white = random.normalvariate(0.0, 1.0)
    last_blue = 0.0
    
    while True:
        # 1. Grab next uniform white noise sample
        white = random.normalvariate(0.0, 1.0)
        
        # 2. Apply high-pass difference logic
        # Diffing consecutive white noise points naturally leaves high-frequency shifts
        blue = (white - last_white) + (beta * last_blue)
        
        # 3. Cache historical registers
        last_white = white
        last_blue = blue
        
        yield blue
