# Python Environment

```
$ brew install --cask miniconda
$ conda create -n eternal python=3.10 pytorch torchvision torchaudio -c pytorch
$ conda init zsh

# Close terminal or source shell
$ source ~/.zshrc
$ conda activate eternal
$ conda install -c conda-forge sentence-transformers

# Deactivate eternal env
$ conda deactivate
```