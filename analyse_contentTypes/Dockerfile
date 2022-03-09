FROM ubuntu:16.04

# copy required files
COPY bin/contentAnalyser_linux /ctfanalyser/
COPY report.gohtml /ctfanalyser/
COPY report.gotxt /ctfanalyser/

# copy output folder
COPY task /ctfanalyser/task/

#set working dir
WORKDIR /ctfanalyser

CMD ./contentAnalyser_linux -d task/ ; cp *.html task/ ; rm *.html

