// Command adminctl manages staff accounts from the server shell.
//
// Admin access is deliberately NOT self-service: there is no public "become an
// admin" page and no admin bootstrap endpoint on the web surface, because an
// endpoint that can mint an administrator is an endpoint an attacker can try.
// Instead an operator runs this tool where the database already trusts them —
// over SSH, or `docker compose exec` — the same shape as Django's
// createsuperuser or Rails' rails console.
//
//	DATABASE_URL=postgres://... adminctl create  -email you@example.com -first Имя -last Фамилия
//	DATABASE_URL=postgres://... adminctl promote -email you@example.com [-role admin]
//	DATABASE_URL=postgres://... adminctl list
//
// The password is never taken from a flag (flags leak into shell history and
// `ps`): it is read from the terminal without echo, or from the ADMIN_PASSWORD
// environment variable for unattended provisioning.
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
	"shanraq.org/pkg/modules/auth"
)

const usage = `adminctl — staff account management

  adminctl create  -email <e-mail> -first <Имя> -last <Фамилия> [-middle <Отчество>] [-role admin]
  adminctl promote -email <e-mail> [-role admin]
  adminctl list

Environment:
  DATABASE_URL    required, PostgreSQL DSN
  ADMIN_PASSWORD  optional, used by "create" instead of the interactive prompt
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fail("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fail("connect: %v", err)
	}
	defer pool.Close()
	store := auth.NewStore(pool)

	switch os.Args[1] {
	case "create":
		cmdCreate(ctx, store, os.Args[2:])
	case "promote":
		cmdPromote(ctx, store, os.Args[2:])
	case "list":
		cmdList(ctx, pool)
	default:
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
}

func cmdCreate(ctx context.Context, store *auth.Store, args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	email := fs.String("email", "", "e-mail of the new administrator")
	first := fs.String("first", "", "given name")
	last := fs.String("last", "", "family name")
	middle := fs.String("middle", "", "patronymic (optional)")
	role := fs.String("role", "admin", "role: admin | director")
	_ = fs.Parse(args)

	normEmail, ok := auth.NormalizeEmail(*email)
	if !ok {
		fail("invalid e-mail")
	}
	f, l, m := auth.NormalizePersonName(*first), auth.NormalizePersonName(*last), auth.NormalizePersonName(*middle)
	if err := auth.ValidatePersonName(f); err != nil {
		fail("first name: %v", err)
	}
	if err := auth.ValidatePersonName(l); err != nil {
		fail("last name: %v", err)
	}
	if err := auth.ValidateOptionalPersonName(m); err != nil {
		fail("patronymic: %v", err)
	}
	if !isStaffRole(*role) {
		fail("role must be admin or director")
	}

	password, err := readPassword()
	if err != nil {
		fail("%v", err)
	}
	if err := auth.ValidatePassword(password); err != nil {
		fail("password: %v", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fail("hash password: %v", err)
	}
	user, err := store.CreateUserNamed(ctx, normEmail, string(hash), f, l, m, *role)
	if err != nil {
		if errors.Is(err, auth.ErrEmailExists) {
			fail("that e-mail already has an account — use `adminctl promote` instead")
		}
		fail("create: %v", err)
	}
	fmt.Printf("created %s (%s) with role %s\n", normEmail, user.ID, *role)
}

func cmdPromote(ctx context.Context, store *auth.Store, args []string) {
	fs := flag.NewFlagSet("promote", flag.ExitOnError)
	email := fs.String("email", "", "e-mail of an existing account")
	role := fs.String("role", "admin", "role: admin | director")
	_ = fs.Parse(args)

	normEmail, ok := auth.NormalizeEmail(*email)
	if !ok {
		fail("invalid e-mail")
	}
	if !isStaffRole(*role) {
		fail("role must be admin or director")
	}
	found, err := store.SetPrimaryRole(ctx, normEmail, *role)
	if err != nil {
		fail("promote: %v", err)
	}
	if !found {
		fail("no account with that e-mail")
	}
	fmt.Printf("%s is now %s\n", normEmail, *role)
}

func cmdList(ctx context.Context, pool *pgxpool.Pool) {
	rows, err := pool.Query(ctx, `
		SELECT email, role, COALESCE(first_name,''), COALESCE(last_name,''), created_at::date
		FROM auth_users
		WHERE role IN ('admin','director','manager','editor')
		ORDER BY created_at`)
	if err != nil {
		fail("list: %v", err)
	}
	defer rows.Close()
	n := 0
	for rows.Next() {
		var email, role, first, last string
		var created time.Time
		if err := rows.Scan(&email, &role, &first, &last, &created); err != nil {
			fail("scan: %v", err)
		}
		fmt.Printf("%-32s %-10s %s %s  (%s)\n", email, role, first, last, created.Format("2006-01-02"))
		n++
	}
	if n == 0 {
		fmt.Println("no staff accounts — create one with `adminctl create`")
	}
}

func isStaffRole(r string) bool { return r == "admin" || r == "director" }

// readPassword takes the password from ADMIN_PASSWORD when set (unattended
// provisioning), otherwise prompts twice without echo.
func readPassword() (string, error) {
	if pw := os.Getenv("ADMIN_PASSWORD"); pw != "" {
		return pw, nil
	}
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		// Piped input: read a single line.
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil && line == "" {
			return "", fmt.Errorf("read password: %w", err)
		}
		return strings.TrimRight(line, "\r\n"), nil
	}
	fmt.Print("Password: ")
	first, err := term.ReadPassword(fd)
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("read password: %w", err)
	}
	fmt.Print("Repeat:   ")
	again, err := term.ReadPassword(fd)
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("read password: %w", err)
	}
	if string(first) != string(again) {
		return "", errors.New("passwords do not match")
	}
	return string(first), nil
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "adminctl: "+format+"\n", args...)
	os.Exit(1)
}
