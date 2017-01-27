# Bounded Networks Problem: A Tree-Recursive Solution

## Problem Statement

You are given `N` networks and you need to return `M` networks where the `M` networks contain the IPs from the `N` but minimise the number of additional IPs that are not in the `M` networks.

## Definitions

   * An *address* is an integer in set `A = [0, 2^32)`.
   * A *network* is a subset of `A` (i.e. a set of addresses) of the form `[ a * 2^k, (a + 1) * 2^k )` for some `a` and `k`.
   * The *footprint* of a sequence of networks is the union of the networks.
   * The *footprint size* of a sequence of networks (or just *size* of a sequence of networks) is the cardinality of the footprint, i.e. the number of addresses in the union of the networks.
   
## Initial Observations

*Left and Right Subnetworks*: for any `[a, b)` of size greater than 1 (i.e. `b > a+1`), its addresses can be partitioned into a *Left Subnetwork* `[a, (a+b)/2)` and a *Right Subnetwork* `[(a+b)/2, y)`.

## Solution Overview

We assume a sequence of input networks `Input := (p[0], ..., p[N-1])`.

For any set `x` of elements of `Input`

   * define `Presolutions(M, x)` to be the set of *presolutions* where a presolution is a sequence of at most `M` networks whose footprint is a superset of the footprint of `x`
   * define `MinSize(M, x)` to be the smallest footprint size of all presolutions in `Presolutions(M, x)` which is `min{ size(union(y)) for y in Presolutions(M, x) }`
   * define `Solutions(M, x)` to be the subset of `Presolutions(M, x)` whose elements have footprint size `MinSize(M, x)`.
   
For non-empty `x`, `Solutions(1, x)` is guaranteed to be a singleton set, so we can

   * define `LeastNetwork(x)` to be the element in the singleton set `Solutions(1, x)`.

Also, we clearly have

    MinSize(M, Input) == MinSize(M, (LeastNetwork(Input),))
    
and

    Solutions(M, Input) == Solutions(M, (LeastNetwork(Input},))

For network `q` in `Input`, we define subsets of `Input`
   * `Left(q)` to be the set of networks which are subsets of the Left Subnetwork of `q`,
   * `Right(q)` to be the set of networks which are the subsets of the Right Subnetwork of `q`.

The solution is based on the observation that for `y = LeastNetwork(x)`, we will find at least one element of `Solutions(M, y)` by considering all presolutions constructed by the concatenation of solutions in `Solutions(j, Left(y))` and `Solutions(M-j, Right(y))` for `j=1..(M-1)`.
   
This observation gives us the following recursive formulation of `MinSize(M, x)`:

    MinSize(M, x) == min{ MinSize(j, Left(LeastNetwork(x))) + MinSize(M-j, Right(LeastNetwork(x))) for 1 <= j <= M-1 }

## Solution

To compute `MinSize(M, Input)` efficiently, we construct a binary tree of networks as follows:

   * `LeastNetwork(Input)` is the root node of the tree
   * for each node `q` if `q` is in `Input`, then it is a leaf
   * otherwise `q` has child nodes `LeastNetwork(Left(q))` and `LeastNetwork(Right(q))`
   
Note that for any nonempty subset `x` of `Input`, both `Left(LeastNetwork(x))` and `Right(LeastNetwork(x))` are non-empty.  To see why this is so, if `Left(LeastNetwork(x))` is empty then `Right(LeastNetwork(x))` is a superset of `x` and strictly smaller than `LeastNetwork(x)` which is a contradiction.  The same argument applies if `Right(LeastNetwork(x))` is empty.

Assuming no overlapping networks in `Input`, the tree has `N` leaves and since it is a binary tree, it must have `N-1` inner nodes.

This tree can be computed efficiently if `Input` is sorted and overlapping networks are removed.  ll subsequences of `Input` of interest can be represented as pairs of indices. `LeastNetwork(x)` can be computed in `O(1)` time for any `x`.  The subsets `Left(q)` and `Right(q)` can determined in `log(size(q))` via binary search.  Since there are `O(N)` tree nodes, the tree can be computed in `O(N * log(N))` time (including sorting and overlap removal).

We compute `MinSize(j, q)` for each tree node `q` (children first) for `1<=j<=K` where `K` is `M` less the distance from the root node.

We can find a solution in `Solutions(M, LeastNetwork(Input))` by traversing the tree again and picking a `j` which attains the minimum at each node.  Of course, this `j` can be selected and stored while the minimum is computed to make the solution finding phase more efficient.  Starting from the root node `q` and bound `M`,

   * if `M==1`, emit `LeastNetwork(q)`
   * if `j>1`, do not emit any network but traverse the left child with bound `j` and the right child with bound `M-j`.
   
### Asymptotic Complexity

   * Sorting input networks: `O(N * log(N))`
   * Removing overlapping networks: `O(N)`
   * Building tree of networks: `O(N * log(N))`
   * Computing `MinSize(M, q)`: `O(M * N)`
   * Find a solution using the tree: `O(N)`

Overall: `O(N * log(N)) + O(M * N)`