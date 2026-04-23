package interview

import (
	"testing"

	"github.com/Jaxongir1006/ai-interview-prep-api/tests/state/database"
)

func GivenDefaultCatalog(t *testing.T) {
	t.Helper()

	db := database.GetTestDB(t)
	ctx, cancel := database.QueryContext()
	defer cancel()

	_, err := db.ExecContext(ctx, `
INSERT INTO interview.interview_target_roles (key, name, description, display_order)
VALUES
    ('python', 'Python Developer', 'Backend and application interviews focused on Python', 10),
    ('golang', 'Go Developer', 'Backend interviews focused on Go services', 20),
    ('javascript', 'JavaScript Developer', 'Frontend and full-stack interviews focused on JavaScript', 30),
    ('java', 'Java Developer', 'Backend interviews focused on Java services', 40);

INSERT INTO interview.interview_experience_levels (key, name, description, display_order)
VALUES
    ('junior', 'Junior', 'Early-career interview expectations', 10),
    ('mid', 'Mid-level', 'Intermediate interview expectations', 20),
    ('senior', 'Senior', 'Senior interview expectations', 30);

INSERT INTO interview.interview_topics (key, name, description, category, display_order)
VALUES
    ('algorithms', 'Algorithms', 'Data structures, complexity, and problem-solving', 'core', 10),
    ('system-design', 'System Design', 'Scalable architecture and distributed systems', 'backend', 20),
    ('database-design', 'Database Design', 'Relational modeling, indexing, and query design', 'backend', 30);

INSERT INTO interview.interview_target_role_topics (target_role_id, topic_id, display_order)
SELECT r.id, t.id, t.display_order
FROM interview.interview_target_roles r
CROSS JOIN interview.interview_topics t
WHERE r.key IN ('python', 'golang', 'javascript', 'java')
  AND t.key IN ('algorithms', 'system-design', 'database-design');
`)
	if err != nil {
		t.Fatalf("GivenDefaultCatalog: failed to create catalog: %v", err)
	}
}
