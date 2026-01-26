# Merging Modular Requirement Files

The `srequirements merge` command provides a powerful way to manage large, complex requirements documents (like PRDs, MRDs, etc.) by breaking them into smaller, more manageable modular files. This document explains when and how to use this feature.

## When to Use `merge`

Managing a single, large JSON or YAML file for a complete requirements document can be cumbersome. It often leads to merge conflicts, difficulty in collaboration, and a lack of clarity on where to update specific sections.

The `merge` command is ideal when you want to:

-   **Promote Collaboration**: Different teams or individuals can work on separate sections of the document simultaneously (e.g., Product Managers on user stories, Engineers on architecture).
-   **Improve Readability**: Smaller, focused files are easier to read and maintain.
-   **Reduce Merge Conflicts**: Changes are isolated to specific files, minimizing the chance of conflicting edits in a monolithic file.
-   **Enable Reusability**: Common sections, like personas, can potentially be reused across different documents.

The workflow involves two key activities:
1.  **Splitting**: A large document is broken down into logical, self-contained JSON files.
2.  **Merging**: The modular files are combined back into a single, complete document for validation, rendering, or distribution.

## How to Use `merge`

The command recursively merges two or more JSON files into a single output file. The merge is "deep," meaning nested objects are merged intelligently. If the same key exists in multiple files, the value from the *last* file specified in the command will take precedence.

### Syntax

```bash
srequirements merge -o <output-file.json> <input-file-1.json> <input-file-2.json> ...
```

-   `-o, --output <file>`: Specifies the path for the final merged JSON file.
-   `<input-file-n.json>`: A space-separated list of the modular JSON files to merge. The order is important.

### Example

Using the example files located in the `examples/merge/` directory, we can construct a complete Product Requirements Document.

The modular files are:
- `prd-metadata.json`
- `prd-objectives.json`
- `prd-personas.json`
- `prd-user-stories.json`
- `prd-requirements.json`
- `prd-roadmap.json`
- `prd-architecture.json`

To merge them into a single `prd-complete.json`, run the following command:

```bash
srequirements merge -o prd-complete.json \
  examples/merge/prd-metadata.json \
  examples/merge/prd-objectives.json \
  examples/merge/prd-personas.json \
  examples/merge/prd-user-stories.json \
  examples/merge/prd-requirements.json \
  examples/merge/prd-roadmap.json \
  examples/merge/prd-architecture.json
```

This will create `prd-complete.json` in your current directory, containing the contents of all seven files.

---

## Instructions for AI Assistants

As an AI assistant, you can use this `merge` functionality to help users manage large structured documents.

### Task: Splitting a Large File

If a user provides a large, monolithic JSON requirements file, you can break it down into logical modules.

**Workflow:**
1.  **Analyze the root keys** of the JSON document (e.g., `metadata`, `personas`, `roadmap`).
2.  **For each root key**, create a new file named `prd-<key-name>.json`.
3.  **Extract the corresponding object** and write it to the new file.
4.  **Inform the user** that you have split the file and list the new modular files you created.

**Example Prompt:** "This PRD file is getting too large. Can you split it into modular files for me?"

### Task: Merging Modular Files

If a user asks you to assemble a complete document from several modular files, you should use the `merge` command.

**Workflow:**
1.  **Identify the modular files** the user wants to merge. Ask for clarification if needed.
2.  **Determine the correct merge order.** A common order is Metadata -> Objectives -> Personas -> Stories -> Requirements -> Roadmap -> Architecture. If unsure, list the files and ask the user to confirm the order.
3.  **Construct and execute the `srequirements merge` command.** Be sure to specify a clear and descriptive output file name.
4.  **Confirm the action** by stating that you have merged the files into the new output file.

**Example Prompt:** "Please merge our PRD files into a single document for final review."
