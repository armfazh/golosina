## Inline vs Call Experiment

*Purpose*: This branch shows an experiment on the effects of inlining code versus calling Go functions.

#### Setup
Clone this repo and change to the branch inlineVScall

````
 $ git clone --branch inlineVScall --single-branch https://github.com/armfazh/golosina
````

To run the code with calls:
````
 $ go test -v -bench=. -benchmem -cover -count=11 -tags=call | tee call.txt
````

To run the inlined code:
````
 $ go test -v -bench=. -benchmem -cover -count=11 -tags=inline | tee inline.txt
````

#### Time

The benchmark was preformed in a Core i7-8650U processor giving these results:

**call(old) vs inlined (new)**

````
name              old time/op    new time/op    delta
ScalarBaseMult-8    23.8µs ± 0%    19.5µs ± 1%  -17.94%  (p=0.000 n=10+10)
ScalarMult-8        39.7µs ± 0%    30.0µs ± 0%  -24.44%  (p=0.000 n=10+11)

name              old speed      new speed      delta
ScalarBaseMult-8  1.34MB/s ± 0%  1.64MB/s ± 1%  +22.16%  (p=0.000 n=11+10)
ScalarMult-8       810kB/s ± 0%  1070kB/s ± 0%  +32.10%  (p=0.000 n=11+10)

name              old alloc/op   new alloc/op   delta
ScalarBaseMult-8     0.00B          0.00B          ~     (all equal)
ScalarMult-8         0.00B          0.00B          ~     (all equal)

name              old allocs/op  new allocs/op  delta
ScalarBaseMult-8      0.00           0.00          ~     (all equal)
ScalarMult-8          0.00           0.00          ~     (all equal)
````
#### Size

**call(old) vs inlined (new)**

| Version | Operation | Size (bytes) | Factor |
|---------|----------------|--------:|-------:|
| Call    | ScalarBaseMult |   1,246 |  1.0x  |
| Inline  | ScalarBaseMult |   5,937 |  4.6x  |
| Call    | ScalarMult     |   1,407 |  1.0x  |
| Inline  | ScalarMult     |  11,780 |  8.2x  |

#### Comments

For this specific workload, inlining code can offer maximum speedup saving around 20% of the time, however it also increases code size significantly.


Authored by Armando Faz [(@armfazh)](https://www.github.com/armfazh)
