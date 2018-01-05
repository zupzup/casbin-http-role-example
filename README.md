# casbin-http-role-exampe

Simplistic Example of role-based HTTP Authorization with [casbin](https://github.com/casbin/casbin) using [scs](https://github.com/alexedwards/scs) for session handling.

Run with

```bash
dep ensure
go run main.go
```

Which starts a server at `http://localhost:8080` with the following routes:

* `POST /login` - accessible if not logged in
   * takes `name` as a form-data parameter - there is no password
   * Valid Users: 
     * `Admin` ID: `1`, Role: `admin`
     * `Sabine` ID: `2`, Role: `member`
     * `Sepp` ID: `3`, Role: `member`
* `POST /logout` - accessible if logged in
* `GET /member/current` - accessible if logged in as a member
* `GET /member/role` - accessible if logged in as a member
* `GET /admin/stuff` - accessible if logged in as an admin
