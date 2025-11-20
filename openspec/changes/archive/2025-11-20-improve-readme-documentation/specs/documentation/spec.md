## ADDED Requirements

### Requirement: Project Directory Structure Documentation
The README SHALL provide a comprehensive directory structure explanation for the full-stack Go+React project.

#### Scenario: Developer understands project layout
- **WHEN** a new developer reads the directory structure section
- **THEN** they understand the separation between frontend (ui/) and backend (root level) code

#### Scenario: Build process comprehension
- **WHEN** reviewing the directory structure documentation
- **THEN** it explains how frontend builds integrate with backend static file serving

### Requirement: Local Development Workflow Documentation  
The README SHALL provide complete step-by-step instructions for local development setup and workflow.

#### Scenario: Frontend development setup
- **WHEN** a developer wants to work on the React frontend
- **THEN** the documentation provides clear instructions for UI development and hot-reload

#### Scenario: Backend development setup
- **WHEN** a developer wants to work on the Go backend
- **THEN** the documentation explains how to run the server and configure dependencies

#### Scenario: Full-stack development workflow
- **WHEN** a developer needs to work on both frontend and backend
- **THEN** the documentation explains the complete build process including frontend compilation and backend integration

### Requirement: Build Process Documentation
The README SHALL explain the relationship between development scripts and production deployment.

#### Scenario: Build script understanding
- **WHEN** a developer reviews the build documentation
- **THEN** they understand the purpose of buildui.sh, build.sh, and Dockerfile relationships

#### Scenario: Production deployment preparation
- **WHEN** preparing for production deployment
- **THEN** the documentation clearly explains the multi-stage build process from development to Docker container