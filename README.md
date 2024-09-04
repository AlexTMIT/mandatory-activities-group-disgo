# The Dining Philosophers in Go

Welcome to the **DisGo** group's GitHub repository for the mandatory activity in Week 2 of the Distributed Systems course at ITU University 2024. This project is part of the 3rd-semester curriculum and focuses on solving the classic concurrency problem known as **The Dining Philosophers**.

## Project Overview

The **Dining Philosophers** problem is a fundamental example in computer science that demonstrates the challenges of concurrency and synchronization. The problem is set at a round dining table with five philosophers who spend their time alternately **eating** and **thinking**.

Each philosopher needs two forks to eat, but there are only five forks, one between each pair of neighboring philosophers. Therefore, at most two philosophers can eat at the same time. The challenge is to design a system that prevents a deadlock situation where no philosopher can eat.

## Objectives

The goal of this project is to implement the dining philosophers problem in Go, with the following requirements:

- **Concurrency with Goroutines**: Each fork and each philosopher must run in its own thread (goroutine).
- **Channel Communication**: Philosophers and forks must communicate solely using Go channels.
- **Deadlock-Free Design**: The system must be designed to prevent deadlock, ensuring that each philosopher eats at least three times.
- **Asynchronous Requests**: Philosophers should be able to request forks at any time; a sequential approach (executing one philosopher at a time) is not acceptable.
- **State Display**: Philosophers must display any state changes (eating or thinking) during their execution.

## Implementation Details

This implementation uses Go's goroutines and channels to model the philosophers and forks. The solution ensures no deadlock occurs by:

- Assigning a goroutine to each philosopher and fork.
- Utilizing channels for communication between philosophers and forks.
- Ensuring that a philosopher only picks up forks when they are available.
- Avoiding a situation where all philosophers hold one fork and wait indefinitely for the second one.
  
All state changes, such as a philosopher starting to **eat** or **think**, will be printed to the console.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributors

- **DisGo Group Members**:
  - Alex Tilgreen Mogensen <alext@itu.dk>
  - Jakob Sønder <jakso@itu.dk>
  - Sara Ziad Al-Janabi <salj@itu.dk>

## 📞 Contact

For any questions, please contact [alext@itu.dk](mailto:alext@itu.dk).

---
