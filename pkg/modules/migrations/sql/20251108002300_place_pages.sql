-- +goose Up
-- Страницы мест: у каждого места свой адрес, у каждой статьи — своё место.
--
-- Замысел: председатель ЖКХ публикует для своего посёлка, акимат — для своей
-- области, и у этого появляется страница, которую находит поисковик. Тот же
-- справочник geo_nodes уже используется объявлениями, так что место на
-- странице сводит статьи и объявления одного края в одну ленту.
--
-- Адрес берётся не из code: он неоднороден — у Качара это kz-kostanay-kachar,
-- а у Костанайской области g65, и /place/g65 не адрес, а шифр. Slug
-- заполняется из названия транслитерацией при запуске приложения: делать это
-- в SQL нельзя, потому что «ж» переходит в «zh», а translate() умеет только
-- посимвольную замену.
ALTER TABLE geo_nodes ADD COLUMN IF NOT EXISTS slug TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS idx_geo_nodes_slug ON geo_nodes (slug) WHERE slug IS NOT NULL;

-- Место статьи. NULL значит «для всех», и это остаётся значением по умолчанию:
-- статья без места ведёт себя ровно так, как вели себя все статьи до сих пор.
ALTER TABLE articles ADD COLUMN IF NOT EXISTS geo_node_id UUID REFERENCES geo_nodes(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_articles_geo ON articles (geo_node_id) WHERE geo_node_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_articles_geo;
ALTER TABLE articles DROP COLUMN IF EXISTS geo_node_id;
DROP INDEX IF EXISTS idx_geo_nodes_slug;
ALTER TABLE geo_nodes DROP COLUMN IF EXISTS slug;
