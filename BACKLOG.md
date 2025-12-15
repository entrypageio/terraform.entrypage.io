This backlog currently only lists milestone `v0.2` items and they remain open for discussion. The software user API is being revised, so the limitations described in `README.md` are expected to be removed.

## Milestone `v0.2` – more functionality for users

### 10. Expose additional user attributes as read-only fields
- **Status**: To Do
- **Labels**: `type:feature`
- **Milestone**: `v0.2`

**Description**  
Extend `entrypage_io_user` with additional read-only fields based on `userSummaryResponse`: `email`, `picture`, `active`, `identity_provider`.

**Acceptance criteria**
- [ ] Resource schema includes read-only attributes for `email`, `picture`, `active`, `identity_provider`.
- [ ] `Read` populates these fields based on the API response.
- [ ] When the API does not return these fields, they are handled cleanly as `null`/`unknown`.
- [ ] Documentation shows example output with these extra fields.

---

### 11. Add data source `entrypage_io_user`
- **Status**: To Do
- **Labels**: `type:feature`
- **Milestone**: `v0.2`

**Description**  
Add a data source that can look up an existing user via `domain` and `preferred_username` or `id` (`sub`), without Terraform managing the user itself.

**Acceptance criteria**
- [ ] New data source `entrypage_io_user` with input attributes (e.g. `domain` + `preferred_username` or `id`).
- [ ] Uses `FindUser` and/or `FindUserByPreferredUsername`.
- [ ] Returns the same read-only fields as the resource (`id`, `preferred_username`, `email`, `picture`, `active`, `identity_provider`).
- [ ] Behavior when "not found" is clearly defined and documented (error vs empty).

---

### 12. Make HTTP timeout and retry configuration customizable
- **Status**: To Do
- **Labels**: `type:feature`
- **Milestone**: `v0.2`

**Description**  
Make the HTTP timeout of the `apiClient` configurable via provider configuration and optionally add simple retry logic for transient errors.

**Acceptance criteria**
- [ ] Provider schema has, for example, an optional field `http_timeout` (in seconds or as a string).
- [ ] `newClient` uses the configured timeout in the `http.Client`.
- [ ] (Optional) Simple retry on certain 5xx errors with backoff, well documented.
- [ ] Documentation explains how and when to tune this.

---

### 13. Add more examples for common scenarios
- **Status**: To Do
- **Labels**: `type:example`
- **Milestone**: `v0.2`

**Description**  
Add extra example configurations to show realistic use cases for `entrypage_io_user` and future resources/data sources.

**Acceptance criteria**
- [ ] Example for importing an existing user only (without create).
- [ ] Example with multiple users within the same domain.
- [ ] Example that uses a data source (after implementing issue 11).
- [ ] Short explanation per example in README or in a dedicated `examples/README.md`.

---

### 14. Add CONTRIBUTING and maintainer docs
- **Status**: To Do
- **Labels**: `type:docs`
- **Milestone**: `v0.2`

**Description**  
Document how external contributors (and your future self) should work with the repo: how to develop, test, branch, and the PR guidelines.

**Acceptance criteria**
- [ ] `CONTRIBUTING.md` describes how to build the provider locally, run tests, and use the examples.
- [ ] Guidelines for code style, commit messages, and PR review.
- [ ] Reference to CI and the release process.

---

## Milestone `vNext` – JWTs Domain Client API

### 15. Implement JWTs Domain Client API support (`entrypage_io_domain_client` resource)
- **Status**: To Do
- **Labels**: `type:feature`
- **Milestone**: `vNext`

**Description**  
Add support for the "JWTS Domain Client API" of EntryPage so that JWT clients for a domain can be managed via Terraform.

**Acceptance criteria**
- [ ] New `apiClient` functions for: creating a domain client, retrieving (get/list), updating (if applicable), deleting.
- [ ] New resource `entrypage_io_domain_client` with relevant attributes (e.g. `domain`, `client_id`, `client_secret` (sensitive), scopes, etc. – exact schema derived from Swagger).
- [ ] CRUD implementation for `entrypage_io_domain_client` with clean error handling.
- [ ] Documentation with examples showing how to define and use a domain client.
- [ ] (Optional) Data source to look up existing domain clients.

---

### 16. Add acceptance tests using Terraform Plugin Testing Framework
- **Status**: To Do
- **Labels**: `type:tests`
- **Milestone**: `vNext`

**Description**  
Add end-to-end acceptance tests using the official Terraform Plugin Framework test package so we can validate behavior against a real (or simulated) EntryPage API.

**Acceptance criteria**
- [ ] Basic acceptance test for `entrypage_io_user` that performs create/read/delete.
- [ ] Test covers idempotent create when the user already exists.
- [ ] Tests can be run locally against a sandbox/development environment.
- [ ] CI can run these tests conditionally (for example only when an API key is available as a secret).
