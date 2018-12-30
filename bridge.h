#ifndef INCLUDE_BRIDGE_H
#define INCLUDE_BRIDGE_H


#include <stddef.h>
// size_t

#include <stdio.h>
// fdopen
// fclose
// FILE

#include <stdlib.h>
// free

#include <string.h>
// strlen

#include <unistd.h>
// dup
// close

#include <sensors/sensors.h>
#include <sensors/error.h>
// sensors_*
// SENSORS_*


extern void parseErrorWrapper(char const *err, int lineno);
extern void parseErrorWfnWrapper(char const *err, char const *filename, int lineno);
extern void fatalErrorWrapper(char const *proc, char const *err);

typedef char const TCchar;
typedef FILE Tfile;
typedef size_t Tsize;
typedef struct sensors_feature const TCsensors_feature;
typedef struct sensors_subfeature const TCsensors_subfeature;


static inline
void set_sensors_parse_error() {
	sensors_parse_error = &parseErrorWrapper;
}

static inline
void set_sensors_parse_error_wfn() {
	sensors_parse_error_wfn = &parseErrorWfnWrapper;
}

static inline
void set_sensors_fatal_error() {
	sensors_fatal_error = &fatalErrorWrapper;
}


#endif
