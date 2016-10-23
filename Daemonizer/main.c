
/* 
 * File:   main.c
 * Author: girishg
 *
 * Created on 20 October, 2016, 2:05 PM
 */

#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <syslog.h>
#include <sys/wait.h>
#include <string.h>
#include "include/ForkSelf.h"
#define die(err) do { fprintf(stderr, "%s\n", err); exit(EXIT_FAILURE); } while (0);


int main(int argc, char** argv) {

    if (argc < 2) {
        printf("Please provide a path to daemonize\n");
        return (EXIT_FAILURE);
    }

    
    char* progname = argv[1];
    char** args = get_program_args(argv, argc);

    if (args) {
        ForkSelf(progname);
        
        fork_n_exec(progname, args);
    } else {
        printf("no params\n");
        exit(1);
    }


    wait(NULL);
    return (EXIT_SUCCESS);
}

