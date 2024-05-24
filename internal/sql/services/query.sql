-- GetService gets service by name
-- name: GetService :one
SELECT * FROM services WHERE service_name = pggen.arg('service_name');

-- GetAllServices gets all services
-- name: GetAllServices :many
SELECT * FROM services;

-- GetServiceNames gets all service names
-- name: GetServiceNames :many
SELECT service_name FROM services;

-- DeleteService deletes a service
-- name: DeleteService :exec
DELETE FROM services WHERE id=pggen.arg('id');

-- ReportService sets reported to true for a service
-- name: ReportService :exec
UPDATE services SET report_count = report_count + 1 WHERE service_name = pggen.arg('service_name');

-- ResetReportedService resets reported for a service
-- name: ResetReportedService :exec
UPDATE services SET report_count = 0 WHERE service_name = pggen.arg('service_name');

-- PatchServicePassword patches service password
-- name: PatchServicePassword :exec
UPDATE services SET password = pggen.arg('password') WHERE service_name = pggen.arg('service_name');
