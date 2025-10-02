# Shanraq Template Engine
# Шанрак Үлгі Қозғалтқышы / Shanraq Template Engine

## 🚀 Overview / Шолу

The Shanraq Template Engine is a powerful, agglutinative ulgi system designed specifically for the Shanraq framework. It leverages Kazakh language features like morphemes, phonemes, and archetypes (algasqy) to provide high-performance ulgi rendering.

Шанрак Үлгі Қозғалтқышы - бұл Шанрак фреймворгі үшін арнайы жасалған күшті, агглютинативтік үлгі жүйесі. Ол морфема, фонема және архетип сияқты қазақ тілінің ерекшеліктерін пайдаланып, жоғары өнімді үлгі рендерингін қамтамасыз етеді.

## ✨ Features / Мүмкіндіктер

### Core Features / Негізгі Мүмкіндіктер
- **🔄 Agglutinative Syntax** - Uses Kazakh language patterns
- **⚡ High Performance** - Compiled ulgis with caching
- **🧩 Morpheme-based Processing** - Dynamic word composition
- **🎵 Phoneme Optimization** - Sound-based optimizations
- **🏗️ Archetype Patterns (algasqy)** - Reusable ulgi patterns
- **🛠️ Helper Functions** - Rich set of built-in helpers
- **🔧 Filter System** - Data transformation filters
- **📄 Partial Templates** - Reusable ulgi components
- **🎨 Layout System** - Master page layouts

### Shanraq-Specific Features / Шанрак-Арнайы Мүмкіндіктер
- **🔤 Morpheme Engine Integration** - Dynamic word creation
- **🎵 Phoneme Engine Integration** - Sound-based optimizations
- **🏗️ Archetype Engine Integration (algasqy)** - Pattern-based ulgis
- **🇰🇿 Kazakh Language Support** - Native language features
- **⚡ Performance Optimization** - SIMD and caching support

## 📁 File Structure / Файл Құрылымы

```
framework/ulgi/
├── ulgi_qozgaltqys_core.tng      # Core ulgi qozgaltqys
├── ulgi_helpers.tng          # Helper functions
├── ulgi_filters.tng         # Filter functions
├── ulgi_utils.tng           # Utility functions
├── ulgi_example.tng         # Usage example
└── README.md                    # This file
```

## 🚀 Quick Start / Жылдам Бастау

### 1. Initialize Template Engine / Үлгі Қозғалтқышын Инициализациялау

```tenge
// Initialize ulgi qozgaltqys
jasau ulgi_qozgaltqys: TemplateEngine = ulgi_qozgaltqys_jasau();
```

### 2. Prepare Data / Деректерді Дайындау

```tenge
// Prepare ulgi data
jasau data: JsonObject = json_object_create();
json_object_set_string(data, "title", "Shanraq Template Engine");
json_object_set_string(data, "content", "Welcome to Shanraq!");
```

### 3. Render Template / Үлгіні Рендерлеу

```tenge
// Render ulgi
jasau html: jol = ulgi_render(ulgi_qozgaltqys, "home_page", data);
```

## 📝 Template Syntax / Үлгі Синтаксисі

### Variable Interpolation / Айнымалы Интерполяция

```html
<!-- Simple variable -->
<h1>{{ title }}</h1>

<!-- With filters -->
<p>{{ content | upper | trim }}</p>

<!-- Nested object access -->
<p>{{ user.name }}</p>
```

### Conditional Statements / Шартты Мәлімдемелер

```html
<!-- Simple condition -->
{{#eger user.is_admin}}
    <div class="admin-panel">Admin Panel</div>
{{/eger}}

<!-- If-else -->
{{#eger user.is_logged_in}}
    <p>Welcome, {{ user.name }}!</p>
{{#basqa}}
    <p>Please log in.</p>
{{/basqa}}
```

### Loops / Циклдар

```html
<!-- Array iteration -->
{{#each posts}}
    <div class="post">
        <h3>{{ title }}</h3>
        <p>{{ content }}</p>
    </div>
{{/each}}
```

### Helper Functions / Көмекші Функциялар

```html
<!-- Math helpers -->
<p>Total: {{ add price tax }}</p>
<p>Average: {{ divide total count }}</p>

<!-- String helpers -->
<p>{{ concat "Hello" " " "World" }}</p>
<p>{{ join tags ", " }}</p>

<!-- Date helpers -->
<p>Published: {{ format_date post.date "YYYY-MM-DD" }}</p>
```

## 🔧 Helper Functions / Көмекші Функциялар

### Conditional Helpers / Шартты Көмекшілер
- `eger(condition, true_value, false_value)` - If-else logic
- `basqa(condition, true_value, false_value)` - Else logic
- `while(condition, content)` - While loop
- `for(start, end, step, content)` - For loop

### Data Helpers / Деректер Көмекшілері
- `each(array, ulgi)` - Array iteration
- `with(data, ulgi)` - Context setting
- `lookup(obj, key)` - Object property access

### String Helpers / Жол Көмекшілері
- `concat(...args)` - String concatenation
- `join(array, separator)` - Array to string
- `split(text, separator)` - String to array

### Math Helpers / Математика Көмекшілері
- `add(a, b)` - Addition
- `subtract(a, b)` - Subtraction
- `multiply(a, b)` - Multiplication
- `divide(a, b)` - Division

### Shanraq-Specific Helpers / Шанрак-Арнайы Көмекшілер
- `morpheme(word, suffix)` - Morpheme composition
- `phoneme(text)` - Phoneme optimization
- `archetype(name, config)` - Archetype application

## 🎨 Filter System / Сүзгі Жүйесі

### String Filters / Жол Сүзгілері
- `upper` - Convert to uppercase
- `lower` - Convert to lowercase
- `capitalize` - Capitalize first letter
- `title` - Title case
- `trim` - Remove whitespace
- `replace(search, replace)` - String replacement

### Number Filters / Сан Сүзгілері
- `round(decimals)` - Round to decimals
- `ceil` - Round up
- `floor` - Round down
- `abs` - Absolute value

### Shanraq-Specific Filters / Шанрак-Арнайы Сүзгілер
- `morpheme_analyze` - Morpheme analysis
- `phoneme_optimize` - Phoneme optimization
- `archetype_apply` - Archetype application

## 📚 Examples / Мысалдар

### Basic Template / Негізгі Үлгі

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ title }}</title>
</head>
<body>
    <h1>{{ title }}</h1>
    <p>{{ content }}</p>
    
    {{#if user.is_logged_in}}
        <p>Welcome, {{ user.name }}!</p>
    {{/if}}
    
    {{#each posts}}
        <div class="post">
            <h3>{{ title }}</h3>
            <p>{{ content }}</p>
        </div>
    {{/each}}
</body>
</html>
```

### Advanced Template / Жоғары Дәрежелі Үлгі

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ title | title }}</title>
</head>
<body>
    <h1>{{ title | upper }}</h1>
    <p>{{ content | trim | replace "old" "new" }}</p>
    
    {{#if user.is_admin}}
        <div class="admin-panel">
            <h2>Admin Panel</h2>
            <p>Total users: {{ user_count | add 1 }}</p>
        </div>
    {{/if}}
    
    {{#each posts | sort | reverse}}
        <div class="post">
            <h3>{{ title | capitalize }}</h3>
            <p>{{ content | slice 0 100 }}...</p>
            <p>Published: {{ date | format_date "YYYY-MM-DD" }}</p>
        </div>
    {{/each}}
    
    <!-- Shanraq-specific features -->
    <div class="{{ archetype "card" card_config }}">
        <h3>{{ morpheme "jasau" "шы" }}</h3>
        <p>{{ phoneme "optimized_content" }}</p>
    </div>
</body>
</html>
```

## ⚡ Performance Features / Өнімділік Мүмкіндіктері

### Compilation / Компиляция
- **Template Compilation** - Templates are compiled to optimized code
- **Caching** - Compiled ulgis are cached for performance
- **Dependency Tracking** - Automatic dependency management
- **Hot Reloading** - Development-time ulgi updates

### Optimization / Оңтайландыру
- **Morpheme-based Caching** - Cache based on word structure
- **Phoneme Optimization** - Sound-based optimizations
- **Archetype Patterns** - Reusable optimized patterns
- **SIMD Support** - Vectorized operations

## 🔌 Integration / Интеграция

### Server Integration / Сервер Интеграциясы

```tenge
// Initialize ulgi qozgaltqys
jasau ulgi_qozgaltqys: TemplateEngine = ulgi_qozgaltqys_jasau();

// Render ulgi
jasau html: jol = ulgi_render(ulgi_qozgaltqys, "home_page", data);

// Send response
http_response_send(html);
```

### Data Binding / Деректер Байланысы

```tenge
// Prepare ulgi data
jasau data: JsonObject = json_object_create();
json_object_set_string(data, "title", "Shanraq Home");
json_object_set_string(data, "content", "Welcome to Shanraq!");

// Render with data
jasau html: jol = ulgi_render(ulgi_qozgaltqys, "home_page", data);
```

## 🛠️ Development / Дәуелдер

### Template Development / Үлгі Дәуелдері
1. **Create Template** - Write ulgi in `.tng` format
2. **Test Template** - Use development server
3. **Optimize Template** - Apply performance optimizations
4. **Deploy Template** - Production deployment

### Debugging / Дәуелдер
- **Template Errors** - Clear error messages
- **Syntax Highlighting** - IDE support
- **Performance Profiling** - Template performance analysis
- **Memory Debugging** - Memory usage tracking

## 📖 Best Practices / Ең Жақсы Тәжірибелер

### Template Design / Үлгі Дизайны
- **Separation of Concerns** - Keep logic separate from presentation
- **Reusability** - Create reusable components
- **Performance** - Optimize for speed
- **Maintainability** - Keep ulgis clean and organized

### Code Organization / Код Ұйымдастыру
- **Modular Structure** - Organize ulgis logically
- **Naming Conventions** - Use consistent naming
- **Documentation** - Document complex ulgis
- **Version Control** - Track ulgi changes

## 🔍 Troubleshooting / Ақауларды Жою

### Common Issues / Жиі Кездесетін Мәселелер
- **Template Not Found** - Check file paths
- **Syntax Errors** - Validate ulgi syntax
- **Performance Issues** - Check caching and optimization
- **Memory Leaks** - Monitor memory usage

### Debug Tools / Дәуелдер Құралдары
- **Template Validator** - Syntax validation
- **Performance Profiler** - Performance analysis
- **Memory Monitor** - Memory usage tracking
- **Error Logger** - Error tracking and logging

## 📈 Future Roadmap / Болашақ Жол Картасы

### Planned Features / Жоспарланған Мүмкіндіктер
- **Advanced Caching** - More sophisticated caching strategies
- **Template Compilation** - Native code compilation
- **Performance Monitoring** - Real-time performance tracking
- **AI Integration** - Artificial intelligence features

### Community Contributions / Қауымдастық Үлесі
- **Plugin System** - Extensible plugin architecture
- **Theme System** - Template theming
- **Internationalization** - Multi-language support
- **Accessibility** - Accessibility features

## 📚 Documentation / Құжаттама

- [Template Engine Documentation](../qujattama/ulgi-qozgaltqys.md) - Complete documentation
- [API Reference](../qujattama/api/ulgi-api.md) - API reference
- [Examples](../examples/) - Code examples
- [Tutorials](../tutorials/) - Step-by-step tutorials

## 🤝 Contributing / Үлес Қосу

We welcome contributions to the Shanraq Template Engine! Please see our [Contributing Guide](../../CONTRIBUTING.md) for details.

Шанрак Үлгі Қозғалтқышына үлес қосуға қуанамыз! Толық ақпарат үшін [Үлес Қосу Нұсқаулығы](../../CONTRIBUTING.md) қараңыз.

## 📄 License / Лицензия

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.

Бұл жоба MIT Лицензиясы бойынша лицензияланған - толық ақпарат үшін [LICENSE](../../LICENSE) файлын қараңыз.

---

**Conclusion / Қорытынды**: The Shanraq Template Engine provides a powerful, agglutinative ulgi system that leverages Kazakh language features for high-performance web development. With its rich feature set, excellent performance, and easy-to-use syntax, it's the perfect choice for building modern web applications with the Shanraq framework.

**Қорытынды**: Шанрак Үлгі Қозғалтқышы қазақ тілінің ерекшеліктерін пайдаланып, жоғары өнімді веб дәуелдері үшін күшті, агглютинативтік үлгі жүйесін қамтамасыз етеді. Оның бай мүмкіндіктері, тамаша өнімділігі және оңай пайдалану синтаксисімен ол Шанрак фреймворгімен заманауи веб қосымшаларын құру үшін тамаша таңдау.




