/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */

/* 
 * File:   main.c
 * Author: girishg
 *
 * Created on 20 October, 2016, 2:05 PM
 */

#include <stdio.h>
#include <stdlib.h>
#include <syslog.h>
#include <sys/wait.h>
#include <string.h>
#include "include/ForkSelf.h"

/*
 * 
 */
int main(int argc, char** argv) {

    if (argc < 2) {
        printf("Please provide a path to daemonize\n");
        return (EXIT_FAILURE);
    }

    char** args = get_program_args(argv, argc);

    if (args) {
        ForkSelf(argv);
        printf("%s, %s\n", argv[1], args[0]);
        fork_n_exec(argv[1], args);
    } else {
        printf("no params\n");
    }

   wait(NULL);
   return (EXIT_SUCCESS);
}

