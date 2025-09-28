# Go development environment
FROM golang

# Install necessary tools (git for branch info, bash-completion for git prompt path)
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    git \
    bash-completion \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set up Go environment
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

# Create workspace directory
WORKDIR /workspace

# System-wide bash prompt with git branch
# - Source both common locations of git prompt script
# - Configure useful flags (dirty, stash, upstream)
# - Keep PS1 simple and robust (no build-time expansion)
# System-wide Git branch prompt (no expansion during build)
RUN cat <<'EOF' >> /etc/bash.bashrc
# --- Git branch in prompt (system-wide, with colors) ---
[ -f /usr/share/git/git-prompt.sh ] && . /usr/share/git/git-prompt.sh
[ -f /usr/lib/git-core/git-sh-prompt ] && . /usr/lib/git-core/git-sh-prompt
export GIT_PS1_SHOWDIRTYSTATE=1
export GIT_PS1_SHOWSTASHSTATE=1
export GIT_PS1_SHOWUPSTREAM=auto
export GIT_PS1_SHOWCOLORHINTS=1
export PS1="\[\e[32m\]\u@\h \[\e[34m\]\w \[\e[33m\]\$(__git_ps1 '(%s)')\[\e[0m\]\$ "
# --- end ---
EOF
