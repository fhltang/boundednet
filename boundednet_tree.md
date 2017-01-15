# Bounded Networks Problem

## Problem Statement

You are given `N` networks and you need to return `M` networks where the `M` networks contain the IPs from the `N` but minimise the number of additional IPs that are not in the `M` networks.

## Definitions

   * An *address* is an integer in set `A = [0, 2^32)`.
   * A *network* is a subset of `A` (i.e. a set of addresses) of the form `[ a * 2^k, (a + 1) * 2^k )` for some `a` and `k`.
   * The *footprint* of a set of networks is the union of the networks.
   * The *footprint size* of a set of networks (or just *size* of a set of networks) is the cardinality of the footprint, i.e. the number of addresses in the union of the networks.
   
## Initial Observations

*Left and Right Subnetworks*: for any `[a, b)` of size greater than 1 (i.e. `b > a+1`), its addresses can be partitioned into a *Left Subnetwork* `[a, (a+b)/2)` and a *Right Subnetwork* `[(a+b)/2, y)`.

## Solution Overview

We assume a set of input networks `B := {p[0], ..., p[N-1]}`.

For any subset `x` of `B`

   * define `Presolutions(M, x)` to be the set of presolutions where a presolution is a set of at most `M` networks whose footprint is a superset of the footprint of `x`
   * define `MinSize(M, x)` to be `min{ size(union(y)) for y in Presolutions(M, x) }`, i.e. the smallest footprint size of all presolutions in `PreSolution(M, x)`
   * define `Solutions(M, x)` to be the subset of `Presolutions(M, x)` whose elements have footprint size `MinSize(M, x)`.
   
Note that for non-empty `x`, `Solutions(1, x)` is guaranteed to be a singleton set.  We define `LeastNetwork(x)` to be the element in `Solutions(1, x)`.

Also, we clearly have

    MinSize(M, B) == MinSize(M, {LeastNetwork(B)})
    
and

    Solutions(M, B) == Solutions(M, {LeastNetwork(B)})

For network `q`, we define subsets of `B`
   * `Left(q)` to be the set of networks which are subsets of the Left Subnetwork of `q`,
   * `Right(q)` to be the set of networks which are the subsets of the Right Subnetwork of `q`.
   
The solution is based on the following recursive formulation of `MinSize(M, x)`:

    MinSize(M, x) == min{ MinSize(j, Left(LeastNetwork(x))) + MinSize(M-j, Right(LeastNetwork(x))) for 1 <= j <= M-1 }

For any set `x` of networks from `B`, we can find the least network `LeastNetwork(x)` which is a superset of the footprint of `x`.  This least network partitions `x` into Left and Right subsets.  For each `j` between `1` and `M-1`, we can construct a presolution by taking a solution of up to size `j` from the left partition together with a solution of up to size `M-j` from the right partition.  Clearly there is a solution in `Solutions(M, x)` amongst one of these presolutions.

## Solution

To compute `MinSize(M, B)` efficiently, we construct a binary tree of networks as follows:

   * `LeastNetwork(B)` is the root node of the tree
   * for each node `q` if `q` is in `B`, then it is a leaf
   * otherwise `q` has child nodes `LeastNetwork(Left(q))` and `LeastNetwork(Right(q))`
   
Note that for any nonempty subset `x` of `B`, both `Left(LeastNetwork(x))` and `Right(LeastNetwork(x))` are non-empty.  To see why this is so, if `Left(LeastNetwork(x))` is empty then `Right(LeastNetwork(x))` is a superset of `x` and strictly smaller than `LeastNetwork(x)` which is a contradiction.  The same argument applies if `Right(LeastNetwork(x))` is empty.

Assuming no overlapping networks in `B`, the tree has `N` leaves and since it is a binary tree, it must have `N-1` inner nodes.

This tree can be computed efficiently if `B` is sorted and overlapping networks are removed.  By enumerating the remaining non-overlapping networks in `B` as `p[0], ..., p[N-1]`, all subsets of `B` of interest can be represented as pairs of indices. `LeastNetwork(x)` can be computed in `O(1)` time for any `x`.  The subsets `Left(q)` and `Right(q)` can determined in `log(size(q))` via binary search.  Since there are `O(N)` tree nodes, the tree can be computed in `O(N * log(N))` time (including sorting and overlap removal).

We compute `MinSize(j, q)` for each node (children first) for `1<=j<=K` where `K` is `M` less the distance from the root node.

We can find a solution in `Solutions(M, LeastNetwork(B))` by traversing the tree again and picking a `j` which attains the minimum at each node.  Starting from the root node `q` and bound `M`,

   * if `M==1`, emit `LeastNetwork(q)`
   * if `j>1`, do not emit any network but traverse the left child with bound `j` and the right child with bound `M-j`.
   
### Asymptotic Complexity

   * Sorting input networks: `O(N * log(N))`
   * Removing overlapping networks: `O(N)`
   * Building tree of networks: `O(N * log(N))`
   * Computing `MinSize(M, q)`: `O(M * N)`
   * Find a solution using the tree: `O(N)`

Overall: `O(N * log(N)) + O(M * N)`