"""The report must not treat source-controlled text as HTML."""

from sholu.bet import _cell


def test_text_is_escaped_in_an_isolated_file(tmp_path):
    page = tmp_path / "fragment.html"
    page.write_text(_cell("аты", "<script>alert(1)</script>"), encoding="utf-8")
    rendered = page.read_text(encoding="utf-8")
    assert "<script>" not in rendered
    assert "&lt;script&gt;" in rendered
