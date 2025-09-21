# Go development environment
FROM golang

# Install necessary tools (git for branch info, bash-completion for git prompt path, locales for proper color handling)
RUN apt-get update && apt-get install -y \
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

# Configure colorful bash prompt with current git branch
# Use __git_ps1 from git-prompt.sh if available; fallback to simple prompt
RUN echo 'if [ -f /usr/lib/git-core/git-sh-prompt ]; then \n\t. /usr/lib/git-core/git-sh-prompt; \nfi' >> /root/.bashrc \
    && echo 'export GIT_PS1_SHOWDIRTYSTATE=1' >> /root/.bashrc \
    && echo 'export GIT_PS1_SHOWSTASHSTATE=1' >> /root/.bashrc \
    && echo 'export GIT_PS1_SHOWUPSTREAM=auto' >> /root/.bashrc \
    && echo "export PS1='\[\e[32m\]\u@\h \[\e[34m\]\w\[\e[33m\]$(__git_ps1 " \(%s\)")\[\e[0m\]\n$ '" >> /root/.bashrc \
    && echo 'if ! type __git_ps1 >/dev/null 2>&1; then PS1="\u@\h \w$ "; fi' >> /root/.bashrc
