# Компоненттер жүйесі / Component System

## Шолу / Overview

Shanraq Template Engine-де жаңа компоненттер жүйесі енгізілді. Бұл жүйе шаблондарды модульдік және қайта пайдалануға болатын етіп жасауға мүмкіндік береді.

## Директория құрылымы / Directory Structure

```
betjagy/
├── bolıkter/           # Компоненттер директориясы
│   ├── header.html     # Header компоненті
│   ├── footer.html     # Footer компоненті
│   └── sidebar.html    # Sidebar компоненті (болашақта)
├── better/             # HTML беттер
└── ulgi/               # Шаблон файлдары
```

## Компоненттерді пайдалану / Using Components

### Синтаксис / Syntax

Компоненттерді шаблонда пайдалану үшін мына синтаксисті қолданыңыз:

```html
<!-- Компонентті шақыру / Component call -->
{{> component_name }}
```

### Мысал / Example

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ page.title }}</title>
</head>
<body>
    <!-- Header компоненті -->
    {{> header }}
    
    <!-- Негізгі мазмұн -->
    <main>
        <h1>{{ content.title }}</h1>
        <p>{{ content.description }}</p>
    </main>
    
    <!-- Footer компоненті -->
    {{> footer }}
</body>
</html>
```

## Компоненттерді жасау / Creating Components

### 1. Файл жасау / Create File

Компонент файлын `betjagy/bolıkter/` директориясында жасаңыз:

```bash
# Header компоненті
betjagy/bolıkter/header.html

# Footer компоненті  
betjagy/bolıkter/footer.html

# Sidebar компоненті
betjagy/bolıkter/sidebar.html
```

### 2. Компонент мазмұны / Component Content

Компонент файлында HTML мазмұнын жазыңыз:

```html
<!-- header.html -->
<header class="shanraq-header">
    <div class="container">
        <h1>{{ site.title }}</h1>
        <p>{{ site.subtitle }}</p>
    </div>
</header>
```

```html
<!-- footer.html -->
<footer class="shanraq-footer">
    <div class="container">
        <p>&copy; {{ current_year }} {{ site.name }}</p>
    </div>
</footer>
```

## Жаңа шаблон логикасы / New Template Logic

### template_engine_build_template()

Бұл функция шаблонды компоненттермен құрастырады:

```tenge
atqar template_engine_build_template(template: jol, data: JsonObject) -> jol {
    jasau result: jol = template;
    
    // Компоненттерді табу және ауыстыру
    result = template_engine_process_components(result, data);
    
    // Айнымалыларды ауыстыру
    result = template_engine_replace_variables(result, data);
    
    // Шартты блоктарды өңдеу
    result = template_engine_process_conditionals(result, data);
    
    // Циклдарды өңдеу
    result = template_engine_process_loops(result, data);
    
    // Функцияларды шақыру
    result = template_engine_process_functions(result, data);
    
    qaytar result;
}
```

### template_engine_process_components()

Компоненттерді тауып ауыстырады:

```tenge
atqar template_engine_process_components(template: jol, data: JsonObject) -> jol {
    jasau result: jol = template;
    jasau i: san = 0;
    
    azirshe (i < result.length - 1) {
        eгер (result[i] == '{' && result[i + 1] == '{' && result[i + 2] == '>') {
            // Компонентті табу және ауыстыру
            jasau component_start: san = i;
            jasau component_end: san = template_engine_find_component_end(result, i);
            
            eгер (component_end > component_start) {
                jasau component_call: jol = result.substring(component_start + 3, component_end);
                jasau component_name: jol = string_trim(component_call);
                jasau component_content: jol = template_engine_load_component(component_name);
                
                // Компонентті мазмұнымен ауыстыру
                jasau before: jol = result.substring(0, component_start);
                jasau after: jol = result.substring(component_end + 2);
                result = before + component_content + after;
                
                i = component_start + component_content.length;
            } aitpese {
                i = i + 1;
            }
        } aitpese {
            i = i + 1;
        }
    }
    
    qaytar result;
}
```

## Артықшылықтар / Advantages

### 1. Қайта пайдалану / Reusability
- Компоненттерді бірнеше шаблонда қолдануға болады
- Код дубликациясын болдырмайды

### 2. Оңай дамыту / Easy Development
- Компоненттерді бөлек дамытуға болады
- Жаңартуларды оңай енгізуге болады

### 3. Жылдамдық / Performance
- Компоненттер кэште сақталады
- Рендеринг жылдам орындалады

### 4. Ұйымдастыру / Organization
- Код жақсы ұйымдастырылған
- Компоненттерді оңай табуға болады

## Мысалдар / Examples

### Жай бет / Simple Page

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ title }}</title>
</head>
<body>
    {{> header }}
    
    <main>
        <h1>{{ page.title }}</h1>
        <p>{{ page.content }}</p>
    </main>
    
    {{> footer }}
</body>
</html>
```

### Блог беті / Blog Page

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ post.title }}</title>
</head>
<body>
    {{> header }}
    
    <main>
        <article>
            <h1>{{ post.title }}</h1>
            <p>{{ post.content }}</p>
            <time>{{ post.date }}</time>
        </article>
    </main>
    
    {{> sidebar }}
    {{> footer }}
</body>
</html>
```

## Болашақ жоспарлар / Future Plans

1. **Sidebar компоненті** - Боковая панель
2. **Navigation компоненті** - Навигация
3. **Breadcrumb компоненті** - Хлебные крошки
4. **Alert компоненті** - Уведомления

## Қорытынды / Conclusion

Жаңа компоненттер жүйесі Shanraq Template Engine-ді одан да күшті және икемді етеді. Бұл жүйе арқылы шаблондарды оңай дамытуға және қолдауға болады.



