# Shanraq.org Project Structure

## Root Directory
The root directory contains only essential files for the project:

- **README.md** - Main project documentation
- **package.json** - Node.js project configuration
- **Makefile** - Build scripts and automation
- **Dockerfile** - Containerization configuration
- **index.js** - Application entry point

## Directory Structure

### Core Directories
- **algasqy/** - Archetypes and design patterns
- **artjagy/** - Backend server components
- **betjagy/** - Frontend components (pages, assets, templates)
- **derekter/** - Database layer (ORM, migrations, models)
- **framework/** - Core framework components
- **qurastyru/** - Compiler and transpiler
- **ısker_qisyn/** - Business logic and services

### Documentation (qujattama/)
- **api/** - API documentation
- **architecture/** - System architecture
- **compliance/** - Compliance and regulatory documentation
- **deployment/** - Deployment guides and procedures
- **investors/** - Investor materials and whitepapers
- **reliability/** - Reliability and scalability documentation
- **roadmap/** - Project roadmap and completion reports
- **security/** - Security documentation
- **transaction_core/** - Transaction core documentation
- **user-guide/** - User guides and tutorials

### Configuration (baptaular/)
- **database_baptaular.json** - Database configuration
- **development_baptaular.json** - Development environment configuration
- **package_optimized.json** - Optimized package configuration
- **server_baptaular.json** - Server configuration

### Development Tools (dev_tools/)
- **cli/** - Command-line interface tools
- **formatter/** - Code formatting tools
- **linter/** - Code linting tools
- **sublime/** - Sublime Text support
- **vscode/** - Visual Studio Code support

### Testing (synaqtar/)
- **benchmarks/** - Performance benchmarks
- **demo/** - Demo applications
- **e2e/** - End-to-end tests
- **integration/** - Integration tests
- **unit/** - Unit tests

## File Organization Principles

### 1. Separation of Concerns
- **Core functionality** in dedicated directories
- **Documentation** centralized in `qujattama/`
- **Configuration** centralized in `baptaular/`
- **Development tools** in `dev_tools/`
- **Testing** in `synaqtar/`

### 2. Logical Grouping
- **Related files** grouped together
- **Documentation** organized by topic
- **Configuration** centralized
- **Tools** separated by purpose

### 3. Clean Root Directory
- **Only essential files** in root
- **All documentation** in `qujattama/`
- **All configuration** in `baptaular/`
- **All tools** in `dev_tools/`

## Benefits of This Structure

### 1. Clarity
- **Clear separation** between different types of files
- **Easy navigation** through the project
- **Logical organization** of components

### 2. Maintainability
- **Easy to find** specific types of files
- **Simple to add** new components
- **Clear structure** for new team members

### 3. Scalability
- **Room for growth** in each directory
- **Modular organization** allows independent development
- **Clear boundaries** between different concerns

## Usage Guidelines

### 1. Adding New Files
- **Core functionality** → appropriate core directory
- **Documentation** → `qujattama/` with appropriate subdirectory
- **Configuration** → `baptaular/`
- **Tools** → `dev_tools/`
- **Tests** → `synaqtar/`

### 2. Modifying Existing Files
- **Keep files** in their designated directories
- **Update documentation** when making changes
- **Maintain consistency** with existing structure

### 3. Project Navigation
- **Start with README.md** for project overview
- **Check qujattama/** for documentation
- **Look in core directories** for functionality
- **Use dev_tools/** for development

This structure provides a clean, organized, and maintainable foundation for the Shanraq.org fintech platform.
