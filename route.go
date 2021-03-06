package main

import (
	"context"
	"fmt"

	"dmitri.shuralyov.com/route/github"
	"github.com/shurcooL/users"
)

// shurcoolSeesHomeRouter implements github.Router that
// targets GitHub issues and PRs on home apps for shurcooL user, and
// targets GitHub issues and PRs on github.com for all other users.
type shurcoolSeesHomeRouter struct {
	users users.Service
}

func (r shurcoolSeesHomeRouter) IssueURL(ctx context.Context, owner, repo string, issueID uint64) string {
	if currentUser, err := r.users.GetAuthenticatedSpec(ctx); err == nil && currentUser == shurcool {
		return homeRouter{}.IssueURL(ctx, owner, repo, issueID)
	}
	return github.DotCom{}.IssueURL(ctx, owner, repo, issueID)
}

func (r shurcoolSeesHomeRouter) IssueCommentURL(ctx context.Context, owner, repo string, issueID, commentID uint64) string {
	if currentUser, err := r.users.GetAuthenticatedSpec(ctx); err == nil && currentUser == shurcool {
		return homeRouter{}.IssueCommentURL(ctx, owner, repo, issueID, commentID)
	}
	return github.DotCom{}.IssueCommentURL(ctx, owner, repo, issueID, commentID)
}

func (r shurcoolSeesHomeRouter) PullRequestURL(ctx context.Context, owner, repo string, prID uint64) string {
	if currentUser, err := r.users.GetAuthenticatedSpec(ctx); err == nil && currentUser == shurcool {
		return homeRouter{}.PullRequestURL(ctx, owner, repo, prID)
	}
	return github.DotCom{}.PullRequestURL(ctx, owner, repo, prID)
}

func (r shurcoolSeesHomeRouter) PullRequestCommentURL(ctx context.Context, owner, repo string, prID, commentID uint64) string {
	if currentUser, err := r.users.GetAuthenticatedSpec(ctx); err == nil && currentUser == shurcool {
		return homeRouter{}.PullRequestCommentURL(ctx, owner, repo, prID, commentID)
	}
	return github.DotCom{}.PullRequestCommentURL(ctx, owner, repo, prID, commentID)
}

func (r shurcoolSeesHomeRouter) PullRequestReviewURL(ctx context.Context, owner, repo string, prID, reviewID uint64) string {
	if currentUser, err := r.users.GetAuthenticatedSpec(ctx); err == nil && currentUser == shurcool {
		return homeRouter{}.PullRequestReviewURL(ctx, owner, repo, prID, reviewID)
	}
	return github.DotCom{}.PullRequestReviewURL(ctx, owner, repo, prID, reviewID)
}

func (r shurcoolSeesHomeRouter) PullRequestReviewCommentURL(ctx context.Context, owner, repo string, prID, reviewCommentID uint64) string {
	if currentUser, err := r.users.GetAuthenticatedSpec(ctx); err == nil && currentUser == shurcool {
		return homeRouter{}.PullRequestReviewCommentURL(ctx, owner, repo, prID, reviewCommentID)
	}
	return github.DotCom{}.PullRequestReviewCommentURL(ctx, owner, repo, prID, reviewCommentID)
}

// homeRouter implements github.Router that
// targets GitHub issues on home's issuesapp, and
// targets GitHub pull requests on home's changes app.
//
// THINK: It embeds home, issuesapp, changes app routing details; can it be composed better?
type homeRouter struct{}

func (homeRouter) IssueURL(_ context.Context, owner, repo string, issueID uint64) string {
	return fmt.Sprintf("/issues/github.com/%s/%s/%d", owner, repo, issueID)
}

func (homeRouter) IssueCommentURL(_ context.Context, owner, repo string, issueID, commentID uint64) string {
	return fmt.Sprintf("/issues/github.com/%s/%s/%d#comment-%d", owner, repo, issueID, commentID)
}

func (homeRouter) PullRequestURL(_ context.Context, owner, repo string, prID uint64) string {
	return fmt.Sprintf("/changes/github.com/%s/%s/%d", owner, repo, prID)
}

func (homeRouter) PullRequestCommentURL(_ context.Context, owner, repo string, prID, commentID uint64) string {
	return fmt.Sprintf("/changes/github.com/%s/%s/%d#comment-c%d", owner, repo, prID, commentID)
}

func (homeRouter) PullRequestReviewURL(_ context.Context, owner, repo string, prID, reviewID uint64) string {
	return fmt.Sprintf("/changes/github.com/%s/%s/%d#comment-r%d", owner, repo, prID, reviewID)
}

func (homeRouter) PullRequestReviewCommentURL(_ context.Context, owner, repo string, prID, reviewCommentID uint64) string {
	return fmt.Sprintf("/changes/github.com/%s/%s/%d#comment-rc%d", owner, repo, prID, reviewCommentID)
}
