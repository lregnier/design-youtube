## Context

The README files were written before the hexagonal alignment refactors. The main mismatches are naming (`ResultPublisher` → `EventPublisher`, consumer → subscriber) and the `ProcessingService` method rename, plus the worker's publisher interface collapsing from two methods to one.

## Goals / Non-Goals

**Goals:**
- Keep README diagrams and text in sync with the actual code

**Non-Goals:**
- Restructuring the README format
- Adding new sections
- Updating configuration tables (those are still accurate)
