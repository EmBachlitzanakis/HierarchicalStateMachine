# Hierarchical State Machine Example in Go

This repository demonstrates a simple implementation of a Hierarchical State Machine (HSM) using the `gohsm` library in Go.

## Description

The state machine models a system with three states:

* **StateA:**
    - Initializes a timer set to expire after 2 seconds.
    - Transitions to StateB on `EventB` or to StateC on timeout (`EventTimeout`).
    - Ignores additional `EventA` received in a row.
* **StateB:**
    - Transitions back to StateA on `EventA`.
    - Ignores `EventC`.
* **StateC:**
    - Transitions to StateB on `EventB`.
