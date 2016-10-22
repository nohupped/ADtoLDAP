/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */

/* 
 * File:   ForkSelf.h
 * Author: girishg
 *
 * Created on 20 October, 2016, 4:18 PM
 */

#pragma once
extern void ForkSelf(char** path);
extern void fork_n_exec(char* path, char** args);
char** get_program_args(char** argv, int argc);

