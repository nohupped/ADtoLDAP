
/* 
 * File:   ForkSelf.c
 * Author: girishg
 * 
 * Created on 20 October, 2016, 4:18 PM
 */
#include <stdlib.h>
#include <sys/stat.h>
#include <syslog.h>
#include <sys/types.h>
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
    printf("Continuing with child %d\n", getpid());
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