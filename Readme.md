To compile:

`go build -o password_gen`

To use:

`./password_gen --dict words_light.txt --mode [fast|optimal]`

Note:
To find the optimal solution checking all combinations is required, but it has O(N^4) time complexity.

We can find good, but not globally optimal solution much faster with O(N) complexity. It counts weights of the words and builds map of heaps of words, mapping first letter to the heap (sorted by weight). Then for each word it tries to build a chain of the desired length taking the words from the appropriate heap (which contains words starting from the last letter of the previous word). So it requeres only 2 passes through the original array of words.

First solution is implemented in optimal package, second in fast package.

Assumptions made:
1. English words are only lowercase
2. English words contain only letters