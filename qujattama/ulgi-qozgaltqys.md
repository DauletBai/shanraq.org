# Shanraq Template Engine Documentation
# Шанрак Үлгі Қозғалтқышы Құжаттамасы / Shanraq Template Engine Documentation

## 📋 Overview / Шолу

The Shanraq Template Engine is a powerful, agglutinative ulgi system designed specifically for the Shanraq framework. It leverages Kazakh language features like morphemes, phonemes, and archetypes (algasqy) to provide high-performance ulgi rendering.

Шанрак Үлгі Қозғалтқышы - бұл Шанрак фреймворгі үшін арнайы жасалған күшті, агглютинативтік үлгі жүйесі. Ол морфема, фонема және архетип сияқты қазақ тілінің ерекшеліктерін пайдаланып, жоғары өнімді үлгі рендерингін қамтамасыз етеді.

## 🚀 Features / Мүмкіндіктер

### Core Features / Негізгі Мүмкіндіктер
- **Agglutinative Syntax** - Uses Kazakh language patterns
- **High Performance** - Compiled ulgis with caching
- **Morpheme-based Processing** - Dynamic word composition
- **Phoneme Optimization** - Sound-based optimizations
- **Archetype Patterns (algasqy)** - Reusable ulgi patterns
- **Helper Functions** - Rich set of built-in helpers
- **Filter System** - Data transformation filters
- **Partial Templates** - Reusable ulgi components
- **Layout System** - Master page layouts

### Shanraq-Specific Features / Шанрак-Арнайы Мүмкіндіктер
- **Morpheme Engine Integration** - Dynamic word creation
- **Phoneme Engine Integration** - Sound-based optimizations
- **Archetype Engine Integration (algasqy)** - Pattern-based ulgis
- **Kazakh Language Support** - Native language features
- **Performance Optimization** - SIMD and caching support

## 📝 Template Syntax / Үлгі Синтаксисі

### Basic Syntax / Негізгі Синтаксис

#### Variable Interpolation / Айнымалы Интерполяция
```html
<!-- Simple variable -->
<h1>{{ title }}</h1>

<!-- With filters -->
<p>{{ content | upper | trim }}</p>

<!-- Nested object access -->
<p>{{ user.name }}</p>
```

#### Conditional Statements / Шартты Мәлімдемелер
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

#### Loops / Циклдар
```html
<!-- Array iteration -->
{{#each posts}}
    <div class="post">
        <h3>{{ title }}</h3>
        <p>{{ content }}</p>
    </div>
{{/each}}

<!-- Object iteration -->
{{#each user}}
    <p>{{ key }}: {{ value }}</p>
{{/each}}
```

#### Helper Functions / Көмекші Функциялар
```html
<!-- Math helpers -->
<p>Total: {{ add price tax }}</p>
<p>Average: {{ divide total count }}</p>

<!-- String helpers -->
<p>{{ concat "Hello" " " "World" }}</p>
<p>{{ join tags ", " }}</p>

<!-- Date helpers -->
<p>Published: {{ format_date post.date "YYYY-MM-DD" }}</p>
<p>Age: {{ age user.birth_date }} years</p>
```

### Advanced Syntax / Жоғары Дәрежелі Синтаксис

#### Shanraq-Specific Helpers / Шанрак-Арнайы Көмекшілер
```html
<!-- Morpheme processing -->
<p>{{ morpheme "jasau" "шы" }}</p>

<!-- Phoneme optimization -->
<p>{{ phoneme "optimized_text" }}</p>

<!-- Archetype application -->
<div class="{{ archetype "card" card_config }}">
    <!-- Card content -->
</div>
```

#### Filter Chains / Сүзгі Тізбектері
```html
<!-- Multiple filters -->
<p>{{ content | upper | trim | replace "old" "new" }}</p>

<!-- Custom filters -->
<p>{{ text | morpheme_analyze | phoneme_optimize }}</p>
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

### Date Helpers / Күн Көмекшілері
- `format_date(date, format)` - Date formatting
- `now()` - Current timestamp
- `age(birth_date)` - Age calculation

### Array Helpers / Массив Көмекшілері
- `length(array)` - Array length
- `first(array)` - First element
- `last(array)` - Last element
- `sort(array)` - Sort array
- `reverse(array)` - Reverse array

### Object Helpers / Объект Көмекшілері
- `keys(obj)` - Object keys
- `values(obj)` - Object values
- `has(obj, key)` - Check property existence

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

### Array Filters / Массив Сүзгілері
- `unique` - Remove duplicates
- `sort` - Sort array
- `reverse` - Reverse array
- `slice(start, end)` - Array slice

### Date Filters / Күн Сүзгілері
- `date(format)` - Date formatting
- `time` - Time formatting
- `datetime` - DateTime formatting

### Shanraq-Specific Filters / Шанрак-Арнайы Сүзгілер
- `morpheme_analyze` - Morpheme analysis
- `phoneme_optimize` - Phoneme optimization
- `archetype_apply` - Archetype application

## 🏗️ Template Structure / Үлгі Құрылымы

### File Organization / Файл Ұйымдастыру
```
betjagy/ulgi/                 # Template directory
├── layouts/                  # Layout ulgis
│   ├── main.tng            # Main layout
│   └── admin.tng            # Admin layout
├── partials/                # Partial ulgis
│   ├── header.tng          # Header partial
│   ├── footer.tng          # Footer partial
│   └── sidebar.tng         # Sidebar partial
├── pages/                   # Page ulgis
│   ├── home_page.tng        # Home page
│   ├── blog_page.tng        # Blog page
│   └── contact_page.tng     # Contact page
└── components/              # Component ulgis
    ├── card.tng            # Card component
    ├── button.tng          # Button component
    └── form.tng            # Form component
```

### Template Inheritance / Үлгі Мұрагерлігі
```html
<!-- layout.tng -->
<!DOCTYPE html>
<html>
<head>
    <title>{{ title }}</title>
</head>
<body>
    {{> header }}
    <main>
        {{ content }}
    </main>
    {{> footer }}
</body>
</html>

<!-- page.tng -->
{{#extends "layout"}}
    {{#block "title"}}My Page{{/block}}
    {{#block "content"}}
        <h1>Welcome to my page!</h1>
    {{/block}}
{{/extends}}
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

### Memory Management / Жад Басқаруы
- **Smart Caching** - Intelligent cache invalidation
- **Memory Pooling** - Efficient memory usage
- **Garbage Collection** - Automatic cleanup
- **Resource Management** - Optimized resource usage

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

### Helper Registration / Көмекші Тіркеу
```tenge
// Register custom helper
atqar custom_helper(value: jol) -> jol {
    qaytar "Custom: " + value;
}

json_object_set_function(ulgi_qozgaltqys.helpers, "custom", custom_helper);
```

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

### Testing / Сынау
- **Unit Tests** - Individual ulgi testing
- **Integration Tests** - Full ulgi system testing
- **Performance Tests** - Benchmarking
- **Compatibility Tests** - Cross-browser testing

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

### Performance Optimization / Өнімділік Оңтайландыру
- **Caching** - Use ulgi caching
- **Minimization** - Minimize ulgi size
- **Compilation** - Pre-compile ulgis
- **Monitoring** - Monitor performance metrics

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

---

**Conclusion / Қорытынды**: The Shanraq Template Engine provides a powerful, agglutinative ulgi system that leverages Kazakh language features for high-performance web development. With its rich feature set, excellent performance, and easy-to-use syntax, it's the perfect choice for building modern web applications with the Shanraq framework.




