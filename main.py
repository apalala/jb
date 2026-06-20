#!/usr/bin/env python3
# Copyright (c) 2017-2026 Juancarlo Añez (apalala@gmail.com)
# SPDX-License-Identifier: MIT
from __future__ import annotations

if __name__ == "__main__":
    import sys

    from jb import main as jb_main

    try:
        sys.exit(jb_main())
    except BrokenPipeError:
        sys.exit(0)
