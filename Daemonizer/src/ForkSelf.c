
/* 
 * File:   ForkSelf.c
 * Author: girishg
 * 
 * Created on 20 October, 2016, 4:18 PM
 */
#include <stdlib.h>
#include <sys/stat.h>
#include <syslog.h>
#include <sys/wait.h>
#include <unistd.h>
#include <stdio.h>
#include <string.h>
#include "../include/ForkSelf.h"
#define BUFSIZE 10000
#define die(err) do { fprintf(stderr, "%s\n", err); exit(EXIT_FAILURE); } while (0);
void ForkSelf(char* path) {
    pid_t pid;
    pid = fork();
    
    if (pid < 0) {
        printf("Couldn't fork\n");
        exit(EXIT_FAILURE);
    }

    if (pid > 0) {
        printf("Forked, killing parent %d\n", getpid());
        exit(EXIT_SUCCESS);
    }
    signal(SIGSEGV, handle_seg);

    pid_t forked_pid = getpid();
    char* create_pid_file = get_pid_file();
    if( access( create_pid_file, F_OK ) != -1 ) {
        printf("Stale pid file %s exists. Remove it and re-run\n", create_pid_file);
        exit(1);
    }
    syslog(LOG_NOTICE, "pid file would be %s, %d", create_pid_file, forked_pid);
    FILE *pid_file;
    pid_file = fopen(create_pid_file, "w");
    fprintf(pid_file, "%d\n", (int)forked_pid);
    fclose(pid_file);

    printf("Child %d will be the new parent and session leader\n", getpid());
    umask(0);
    openlog(path, LOG_NOWAIT|LOG_PID,LOG_USER);
    syslog(LOG_NOTICE, "%s Daemonized\n", path);
    
    // New Session ID
    pid_t sid;
    sid = setsid();
    if (sid < 0) {
        syslog(LOG_ERR, "Couldn't create process group\n");
        exit(EXIT_FAILURE);
    }
    // Chdir to var/run
    if ((chdir("/var/run/")) < 0) {
        syslog(LOG_ERR, "Couldn't chdir to /var/run\n");
        exit(EXIT_FAILURE);
    }

}

int pipes[2];
char buf[BUFSIZE];
pid_t child_pid;

void fork_n_exec(char* path, char** args){

    pid_t pid;
    if(pipe(pipes) == -1){
        die("Pipe error");
    }
    pid = fork();  
    switch (pid) {
        case 0:
            printf("new child pid is %d\n", getpid());
            dup2 (pipes[1], STDOUT_FILENO ); // duplicate the handle as STDOUT
            dup2(pipes[1], STDERR_FILENO); // duplicate the handle as STDERR
            close(pipes[0]); // Close both handles.
            close(pipes[1]);
            execv(path, args);
            break;
        default :
            // Closing stdout, stderr and stdin for the parent
            close(STDIN_FILENO);
            close(STDOUT_FILENO);
            close(STDERR_FILENO);
            child_pid = pid;
            signal(SIGHUP, handle_hup_to_sigterm);
            close(pipes[1]); // Closing the writing handle, as we are not using it
            ssize_t nbytes;
            for (;;) {
                nbytes = read(pipes[0], buf, BUFSIZE); // Reading n bytes to buf.
                if (nbytes == -1){
                    break;
                }
                if (nbytes == 0) {
                    break;
                }
         //       write(STDOUT_FILENO, "\n", 1);
                syslog(LOG_NOTICE, "%s", buf);
            }
            wait(NULL);
    }
}

char** get_program_args(char** argv, int argc) {
    if (argc < 3) {
        return NULL;
    }
    return &argv[1];
    
}

void handle_hup_to_sigterm(int signal) {
    syslog(LOG_NOTICE, "Received signal %d, will be killing child pid %d\n", signal, child_pid);
    kill(child_pid, 15);
    
    syslog(LOG_NOTICE, "Killed child pid %d\n", child_pid);
    
}

void handle_seg(int signal) {
    syslog(LOG_ERR, "Received signal %d,(segmentation fault). The program will die now. This could be because of a permission "
            "problem to write the pid file to /var/run, or something serious. Please check if you are "
            "running this program as a privileged user.", signal);
    printf("Received signal %d,(segmentation fault). The program will die now. This could be because of a permission "
                   "problem to write the pid file to /var/run, or something serious. Please check if you are "
                   "running this program as a privileged user.\n", signal);
    exit(1);
}

char* touch_pid;
char*  get_pid_file() {
    char* program = (char*) malloc(strlen(progname) + 1);
    strncat(program, progname, strlen(progname) + 1);
    int i = strlen(program)+1;
    char* filename = (char*) malloc(strlen(progname) + 1);

    while (program[i] != '/') {
        if (i < 0){
            break;
        }
        filename[i] = program[i];
        i--;

    }
    touch_pid = concat("/var/run/", &filename[i+1], ".pid");
    return touch_pid;

}

char* concat(const char *s1, const char *s2, const char *s3)
{
    char *result = malloc(strlen(s1)+strlen(s2)+strlen(s3)+1);
    strncpy(result, s1, strlen(s1));
    strncat(result, s2, strlen(s2));
    strncat(result, s3, (strlen(s3)+1));
    return result;
}