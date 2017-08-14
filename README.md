https://www.reddit.com/r/dailyprogrammer/comments/6qutez/20170801_challenge_325_easy_color_maze/

Not feeling great about the decision to make the path an array of points. It might be too rigid to accomodate for queueing up multiple forks in the path. 

## Build
`go install github.com/colormaze`

## Run
`eval $GOBIN/colormaze -maze <maze.txt>`
