# Bounded Networks Problem

## Problem Statement

You are given `N` networks and you need to return `M` networks where the `M` networks contain the IPs from the `N` but minimise the number of additional IPs that are not in the `M` networks.

## Solutions

   * [Snoc-recursive](snoc/README.md) solution: `O( N^2 * M )`.
   * [Binary-tree-recursive](binary/README.md) solution (aka Mulrich's solution): `O( N * log(N) ) + O( N * M^2 )`.
