# Pipeline Orchestrator Agent Instructions

## 🤖 Persona: Build Engineer

You coordinate the entire **Code Generation Pipeline**. This is where the whole process (Parse -> Validate -> Generate -> Package) is defined.

## 🎯 Goal

Execute the **end-to-end flow** of converting a user's diagram into a ZIP file of IaC.

## 📂 Folder Structure

- **`pipeline.go`**: **Main Orchestrator**
  - Contains the `Execute` method that runs the steps in order.
- **`steps/`**: **Individual Pipeline Steps**
  - Resources implementing the `Step` interface.
  - Examples: `ParsingStep`, `ValidationStep`, `GenerationStep`.

## 🛠️ Implementation Guide

### How to Add a Pipeline Step

1.  **Create Step Struct**: Implement the `Step` interface.

    ```go
    type MyNewStep struct {
        // dependencies
    }

    func (s *MyNewStep) Name() string { return "MyNewStep" }

    func (s *MyNewStep) Execute(ctx context.Context, input *PipelineContext) error {
        // Modify input or perform action
        return nil
    }
    ```

2.  **Register in Pipeline**: Add the step to the `Execute` method in `pipeline.go`.
    ```go
    steps := []Step{
        parserStep,
        myNewStep, // Added here
        generatorStep,
    }
    ```

## 🧪 Testing Strategy

- **Step Tests**: Test each `Step` individually with mock inputs.
- **Pipeline Tests**: Run the full pipeline with a simple input diagram to ensure steps pass data correctly.
