To compile:

`go build -o password_gen`

To use:

`./password_gen --dict words_light.txt --mode [fast|optimal]`

Note:
To find the optimal solution checking all combinations is required, but it has O(N^4) time complexity. This can be speed up if the dictioanry is sorted by weights and we iterate over it frim the start remembering best word combination with desired length. As the atrray is sorted we can stop then the weight starts to grow.

Assumptions made:
1. English words are only lowercase
2. English words contain only letters

We can move any dictionary to the required assumptions, but it can be done with the separate tools.