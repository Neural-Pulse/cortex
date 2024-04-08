### Contributing Guide

#### General

Before contributing, please ensure you have read and understood this guide. We are committed to maintaining a respectful and inclusive collaboration environment. Any contributions not adhering to this code of conduct will be rejected.

#### Frontend (Cortex Frontend)

To contribute to the Cortex project's frontend, please adhere to the following guidelines:

- **TypeScript**: All code should be written in TypeScript. Ensure you follow best practices and typing guidelines to maintain code maintainability and scalability.
- **Conventional Commits**: Use the Conventional Commits format for all commit messages. This aids in automatic changelog generation and maintaining semantic versioning.
- **SemVer**: We adopt SemVer (Semantic Versioning) for project versions. When contributing, consider the changes made and appropriately increment the version following SemVer rules.
- **Code Review**: All Pull Requests must be reviewed by at least one team member before merging. This ensures code quality and consistency.
- **Testing**: Write unit tests, integration tests, and end-to-end (e2e) tests for new features or bug fixes. Testing helps ensure the frontend's stability and functionality.
- **Compile and Build**: Before submitting a pull request, compile and build your changes to ensure everything works as expected.

#### Backend (Cortex Backend)

To contribute to the Cortex project's backend, please adhere to the following guidelines:

- **Go Best Practices**: Follow Go best practices, including effective use of goroutines where appropriate, error handling, and using `context` for request and long-running operation control.
- **Conventional Commits**: Use the Conventional Commits format for all commit messages, as with the frontend.
- **SemVer**: Semantic versioning also applies to the backend. Ensure you version appropriately based on the changes made.
- **Testing**: Write unit tests, integration tests, and ensure your changes do not break existing tests. Test coverage helps ensure the backend's stability.
- **Documentation**: Update documentation as necessary. Code comments and updated READMEs help new contributors better understand the project.

#### How to Contribute

1. **Fork the Repository**: Start by forking the project to your GitHub account.
2. **Clone Your Fork**: Clone your fork to your local machine to start working on the changes.
3. **Create a Branch**: For new features or fixes, create a branch off `main`.
4. **Make Your Changes**: Implement your changes, following the guidelines above.
5. **Commit Your Changes**: Commit your changes using the Conventional Commits format.
6. **Push to Your Fork**: Push your changes to your fork on GitHub.
7. **Open a Pull Request**: In the original repository, open a Pull Request from your branch. Describe the changes made and any other relevant information for the reviewers.
8. **Compile and Build**: Before submitting your pull request, ensure to compile and build your changes to verify everything is working as expected.

Remember to keep your fork synchronized with the original repository to avoid merge conflicts.

We appreciate your contribution to the Cortex project! Together, we can build something amazing.
