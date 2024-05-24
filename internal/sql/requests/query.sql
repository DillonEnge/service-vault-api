-- GetAllRequests gets all requests
-- name: GetAllRequests :many
SELECT * FROM requests;

-- CreateAccessRequest creates a new access request
-- name: CreateAccessRequest :exec
INSERT INTO requests (email, type)
VALUES (pggen.arg('email'), 'access');

-- CreateCodeRequest creates a new 2FA code request
-- name: CreateCodeRequest :exec
INSERT INTO requests (service_name, email, type)
VALUES (pggen.arg('service_name'), pggen.arg('email'), 'code');

-- DeleteRequest deletes a request
-- name: DeleteRequest :exec
DELETE FROM requests WHERE id=pggen.arg('id');
