To compile:

`go build -o password_gen`

To use:

`./password_gen --dict words_light.txt --mode [fast|optimal]`

Note:
To find the optimal solution checking all combinations is required, but it has O(N^4) time complexity.
We can find good, but not globally optimal solution much faster with O(N^2) complexity. 

First solution is implemented in optimal package, second in fast package.

Assumptions made:
1. English words are only lowercase
2. English words contain only letters