# syntax=docker/dockerfile:1
FROM debian:bullseye

# set env
USER ${USER_NAME}
WORKDIR /home/${USER}
SHELL ["/bin/bash", "-c"]
ENV HOME /home/${USER}
ENV PYENV_ROOT $HOME/.pyenv
ENV PATH $PYENV_ROOT/shims:$PYENV_ROOT/bin:$PATH

# install git
RUN apt-get update && apt-get install -y git

# install dependencies
RUN apt-get update && apt-get install -y --no-install-recommends make build-essential libssl-dev zlib1g-dev libbz2-dev libreadline-dev libsqlite3-dev wget curl llvm libncurses5-dev xz-utils tk-dev libxml2-dev libxmlsec1-dev libffi-dev liblzma-dev


# install pyenv
RUN git clone https://github.com/pyenv/pyenv.git .pyenv


# set up bash env && install virtualenv
# install python version and set global
# create & activate virtual env
RUN echo 'export PYENV_ROOT="$HOME/.pyenv"' >> $HOME/.bashrc
RUN echo 'export PATH="$PYENV_ROOT/bin:$PYENV_ROOT/.local/bin:$PATH"' >> $HOME/.bashrc
RUN echo 'eval "$(pyenv init -)"' >> $HOME/.bashrc
RUN pyenv install 3.11 && \
	pyenv global 3.11
RUN git clone https://github.com/pyenv/pyenv-virtualenv.git $PYENV_ROOT/plugins/pyenv-virtualenv 
RUN echo 'eval "$(pyenv virtualenv-init -)"' >> $HOME/.bashrc
RUN pyenv virtualenv 3.11 vote
RUN source $HOME/.bashrc && pyenv activate vote

# build up prefect
# prefect login
RUN pip install -U prefect
RUN --mount=type=secret,id=prefect_env source /run/secrets/prefect_env && prefect cloud login --key ${PREFECT_LOGIN_KEY} --workspace bangvotingguidegmailcom/taiwan-voting-guide 

# install vim
RUN apt-get update && apt-get install -y vim
 
# start the agent
CMD ["/bin/bash", "-c", "prefect agent start -q 'default'"]
