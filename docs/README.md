# Answers

## Task 1

* What happens if you remove the `go-command` from the `Seek` call in the `main` function?

The matching will always be done in order of the array instead of pseudo-randomly when all people seek as goroutines at the same time. I.e. Anna will always match with Bob, Cody will always match with Dave and Eva wont be matched with anyone.

* What happens if you switch the declaration `wg := new(sync.WaitGroup)` to `var wg sync.WaitGroup` and the parameter `wg *sync.WaitGroup` to `wg sync.WaitGroup`?

The WaitGroup gets passed into Seek by value instead of as a pointer, which means Seek only decrements a copy of the WaitGroup in main when calling Done. The result is that main arrives at a deadlock waiting for the WaitGroup which never decrements to 0.

* What happens if you remove the buffer on the channel match?

The person that doesn't get a match gets stuck waiting for another goroutine to receive their message in the match channel as there's no longer a buffer that holds it and lets them continue, and the program therefore arrives at a deadlock.

* What happens if you remove the default-case from the case-statement in the `main` function?

Nothing when there's an odd amount of people as there will always be someone who didn't get matched but if there's an even amount of people, noone is unmatched and the program will arrive at a deadlock waiting for a remaining person in the match channel.

## Task 3

Performances:

|Variant       | Runtime (ms) |
| ------------ | ------------:|
| singleworker |          776 |
| mapreduce    |          430 |