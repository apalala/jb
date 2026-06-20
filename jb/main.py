#!/usr/bin/env python3
# Copyright (c) 2017-2026 Juancarlo Añez (apalala@gmail.com)
# SPDX-License-Identifier: BSD-4-Clause
from __future__ import annotations

import sys


def main() -> int:
    from .jb import main as jb_main

    return jb_main()


if __name__ == "__main__":
    sys.exit(main())
