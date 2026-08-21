-- +goose Up
-- Организация как автор: ЖКХ «Качарец», ТОО «Водоканал», акимат города.
--
-- Организация — не аккаунт, а право аккаунта. За публикацией остаётся человек:
-- у нас правило, что у каждого настоящее имя, чтобы за текстом стоял кто-то
-- конкретный, и общий логин на отдел это правило отменяет — потом не найти,
-- кто писал. Сотрудники меняются, аккаунт живёт дольше человека, а за
-- сообщением акимата по закону стоит ответственное лицо, а не логин.
--
-- Поэтому: в подписи стоит организация, в журнале модерации — Иванов И.И.
--
-- Проверка обязательна и составляет саму суть. Название, которое можно вписать
-- самому, превращает площадку в машину для выдачи себя за акимат, а одна
-- поддельная публикация от его имени стоит дороже, чем польза от сотни
-- настоящих. Пока status не 'verified', название не показывается нигде.
--
-- Форма повторяет re_agents намеренно: там уже работают БИН, очередь модерации
-- и значок проверки, и вторая сущность с теми же потребностями стоила бы
-- дороже, чем общая.
CREATE TABLE IF NOT EXISTS org_authors (
    user_id UUID PRIMARY KEY REFERENCES auth_users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    kind TEXT NOT NULL,
    bin TEXT NOT NULL DEFAULT '',
    -- Территория организации. Акимат Костаная публикует для Костаная и того,
    -- что внутри; иначе проверенный значок становится пропуском на всю страну.
    geo_node_id UUID REFERENCES geo_nodes(id) ON DELETE SET NULL,
    about TEXT NOT NULL DEFAULT '',
    contact TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending',
    reject_reason TEXT NOT NULL DEFAULT '',
    reviewed_at TIMESTAMPTZ,
    reviewed_by UUID REFERENCES auth_users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT org_authors_status_chk CHECK (status IN ('pending','verified','rejected'))
);

CREATE INDEX IF NOT EXISTS idx_org_authors_status ON org_authors (status);
CREATE INDEX IF NOT EXISTS idx_org_authors_place ON org_authors (geo_node_id) WHERE geo_node_id IS NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS org_authors;
