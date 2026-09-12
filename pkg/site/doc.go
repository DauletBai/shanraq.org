// Package site holds what every page of shanraq.org needs regardless of what
// the page is about: the language of the interface, the UI dictionary, the
// shared page context the header and footer read, and the renderer the modules
// put their own templates into.
//
// The rule that keeps this package small: site owns the frame, a module owns a
// subject. Anything that answers "what does this page show" belongs to a module
// under pkg/modules; anything that answers "what does every page have" belongs
// here. A type lives here only when it is the shape a shared partial renders --
// the code that fills it stays with the domain that knows how.
package site
