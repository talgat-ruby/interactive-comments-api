package model

import (
	"context"
)

type CommentReply struct {
	ID          int
	Content     string
	Author      string
	AvatarUrl   string
	Likes       int
	Duration    string
	IsMine      bool
	MyRate      int
	ReplyAuthor string
}

type Comment struct {
	ID        int
	Content   string
	Author    string
	AvatarUrl string
	Likes     int
	Duration  string
	IsMine    bool
	MyRate    int
	Replies   []*CommentReply
}

func (m *Model) GetComments(ctx context.Context, username string) ([]*Comment, error) {
	m.log.InfoContext(ctx, "start GetList")

	comments, err := m.getComments(ctx, username)
	if err != nil {
		m.log.ErrorContext(ctx, "fail GetList", "error", err)
		return nil, err
	}

	comments, err = m.getCommentsReplies(ctx, username, comments)
	if err != nil {
		m.log.ErrorContext(ctx, "fail GetList", "error", err)
		return nil, err
	}

	m.log.InfoContext(ctx, "success GetList")
	return comments, nil
}

func (m *Model) getComments(ctx context.Context, username string) ([]*Comment, error) {
	m.log.InfoContext(ctx, "start getComments")

	sqlStatement := `
		SELECT
			c.OID as id,
			c.content as content,
			u.username as author ,
			u.avatar_url as avatar_url,
			CASE
				WHEN l.count is NULL THEN 0
				ELSE l.count
				END AS likes,
			CASE
				WHEN (strftime('%Y', 'now') - strftime('%Y', c.created_at)) > 0
					THEN 'More than ' || (strftime('%Y', 'now') - strftime('%Y', c.created_at)) || ' year(s) ago'
				WHEN (strftime('%m', 'now') - strftime('%m', c.created_at)) > 0
					THEN 'More than ' || (strftime('%m', 'now') - strftime('%m', c.created_at)) || ' month(s) ago'
				WHEN (strftime('%d', 'now') - strftime('%d', c.created_at)) > 0
					THEN 'More than ' || (strftime('%d', 'now') - strftime('%d', c.created_at)) || ' day(s) ago'
				WHEN (strftime('%H', 'now') - strftime('%H', c.created_at)) > 0
					THEN 'More than ' || (strftime('%H', 'now') - strftime('%H', c.created_at)) || ' hour(s) ago'
				WHEN (strftime('%M', 'now') - strftime('%M', c.created_at)) > 0
					THEN 'More than ' || (strftime('%M', 'now') - strftime('%M', c.created_at)) || ' minute(s) ago'
				WHEN (strftime('%S', 'now') - strftime('%S', c.created_at)) > 0
					THEN 'More than ' || (strftime('%S', 'now') - strftime('%S', c.created_at)) || ' second(s) ago'
				ELSE 'now'
				END AS duration,
			u.username == ? as is_mine,
			CASE
				WHEN l2.rate is NULL THEN 0
				ELSE l2.rate
				END AS my_rate
		FROM
			comment as c
				LEFT JOIN
			main.user_ u
			ON
				c.author = u.username
				LEFT JOIN
			(
				SELECT
					comment_id,
					SUM(rate) as count
				FROM
					like_
				GROUP BY
					comment_id
			) as l
			ON
				c.OID = l.comment_id
				LEFT JOIN
			like_ l2
			ON
				c.OID = l2.comment_id AND l2.author = ?
		WHERE c.parent_id IS NULL
		ORDER BY c.OID;
	`

	rows, err := m.db.Query(sqlStatement, username, username)
	if err != nil {
		m.log.ErrorContext(ctx, "fail getComments", "error", err)
		return nil, err
	}

	defer rows.Close()

	comments := make([]*Comment, 0)

	for rows.Next() {
		c := new(Comment)

		if err = rows.Scan(&c.ID, &c.Content, &c.Author, &c.AvatarUrl, &c.Likes, &c.Duration, &c.IsMine, &c.MyRate); err != nil {
			m.log.ErrorContext(ctx, "fail getComments", "error", err)
			return nil, err
		}

		comments = append(comments, c)
	}

	m.log.InfoContext(ctx, "success getComments")
	return comments, nil
}

func (m *Model) getCommentsReplies(ctx context.Context, username string, cs []*Comment) ([]*Comment, error) {
	m.log.InfoContext(ctx, "start getCommentsReplies")

	sqlStatement := `
		SELECT
			c.id as id,
			c.content as content,
			u.username as author ,
			u.avatar_url as avatar_url,
			CASE
				WHEN l.count is NULL THEN 0
				ELSE l.count
			END AS likes,
			CASE
				WHEN (strftime('%Y', 'now') - strftime('%Y', c.created_at)) > 0
					THEN 'More than ' || (strftime('%Y', 'now') - strftime('%Y', c.created_at)) || ' year(s) ago'
				WHEN (strftime('%m', 'now') - strftime('%m', c.created_at)) > 0
					THEN 'More than ' || (strftime('%m', 'now') - strftime('%m', c.created_at)) || ' month(s) ago'
				WHEN (strftime('%d', 'now') - strftime('%d', c.created_at)) > 0
					THEN 'More than ' || (strftime('%d', 'now') - strftime('%d', c.created_at)) || ' day(s) ago'
				WHEN (strftime('%H', 'now') - strftime('%H', c.created_at)) > 0
					THEN 'More than ' || (strftime('%H', 'now') - strftime('%H', c.created_at)) || ' hour(s) ago'
				WHEN (strftime('%M', 'now') - strftime('%M', c.created_at)) > 0
					THEN 'More than ' || (strftime('%M', 'now') - strftime('%M', c.created_at)) || ' minute(s) ago'
				WHEN (strftime('%S', 'now') - strftime('%S', c.created_at)) > 0
					THEN 'More than ' || (strftime('%S', 'now') - strftime('%S', c.created_at)) || ' second(s) ago'
				ELSE 'now'
			END AS duration,
		    u.username == ? as is_mine,
			CASE
				WHEN l2.rate is NULL THEN 0
				ELSE l2.rate
				END AS my_rate,
		    r.author as reply_author
		FROM
			(
				SELECT
					OID as id,
					reply_id,
					author,
					content,
					created_at,
					parent_id
				FROM comment
				WHERE parent_id == ?
			) as c
				LEFT JOIN
			main.comment r
			ON
				c.reply_id = r.OID
				LEFT JOIN
			main.user_ u
			ON
				c.author = u.username
				LEFT JOIN
			(
				SELECT
					comment_id,
					SUM(rate) as count
				FROM
					like_
				GROUP BY
					comment_id
			) as l
			ON
				c.id = l.comment_id
				LEFT JOIN
			like_ l2
			ON
				c.id = l2.comment_id AND l2.author = ?
		ORDER BY c.id;
	`
	stm, err := m.db.Prepare(sqlStatement)
	if err != nil {
		m.log.ErrorContext(ctx, "fail getCommentsReplies", "error", err)
		return nil, err
	}
	defer stm.Close()

	for _, c := range cs {
		rows, err := stm.Query(username, c.ID, username)
		if err != nil {
			m.log.ErrorContext(ctx, "fail getCommentsReplies", "error", err)
			return nil, err
		}

		list := make([]*CommentReply, 0)

		for rows.Next() {
			c := new(CommentReply)

			if err = rows.Scan(&c.ID, &c.Content, &c.Author, &c.AvatarUrl, &c.Likes, &c.Duration, &c.IsMine, &c.MyRate, &c.ReplyAuthor); err != nil {
				m.log.ErrorContext(ctx, "fail getCommentsReplies", "error", err)
				return nil, err
			}

			list = append(list, c)
		}

		c.Replies = list

		rows.Close()
	}

	m.log.InfoContext(ctx, "success getCommentsReplies")
	return cs, nil
}
