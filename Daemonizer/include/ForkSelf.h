
/* 
 * File:   ForkSelf.h
 * Author: girishg
 *
 * Created on 20 October, 2016, 4:18 PM
 */

#pragma once

#include <sys/types.h>

extern pid_t child_pid;

extern void ForkSelf(char* path);
extern void fork_n_exec(char* path, char** args);
extern char** get_program_args(char** argv, int argc);
extern void handle_hup_to_sigterm(int signal);

