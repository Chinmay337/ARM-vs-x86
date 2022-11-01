
  1
  2
  3
  4
  5
  6
  7
  8
  9
 10
 11
 12
 13
 14
 15
 16
 17
 18
 19
 20
 21
 22
 23
 24
 25
 26
 27
 28
 29
 30
 31
 32
 33
 34
 35
 36
 37
 38
 39
 40
 41
 42
 43
 44
 45
 46
 47
 48
 49
 50
 51
 52
 53
 54
 55
 56
 57
 58
 59
 60
 61
 62
 63
 64
 65
 66
 67
 68
 69
 70
 71
 72
 73
 74
 75
 76
 77
 78
 79
 80
 81
 82
 83
 84
 85
 86
 87
 88
 89
 90
 91
 92
 93
 94
 95
 96
 97
 98
 99
100
101
102
103
104
105
106
107
108
109
110
111
112
113
114
115
116
117
118
119
120
121
122
123
124
125
126
127
128
129
130
131
132
133
134
135
136
137
138
139
140
141
142
143
144
145
146
147
148
149
150
151
152
153
154
155
156
157
158
159
160
161
162
163
164
165
166
167
168
169
170
171
172
173
174
175
176
177
178
179
180
181
182
183
184
185
186
187
188
189
190
191
192
193
194
195
196
197
198
199
200
201
202
203
204
205
206
207
208
209
210
211
212
213
214
215
216
217
218
219
220
221
222
223
224
225
226
227
228
229
230
231
232
233
234
235
236
237
238
239
240
241
242
243
244
245
246
247
248
249
250
251
252
253
254
255
256
257
258
259
260
261
262
263
264
265
266
267
268
269
270
271
272
273
274
275
276
277
278
279
280
281
282
283
284
285
286
287
288
289
290
291
292
293
294
295
296
297
298
299
300
301
302
303
304
305
306
307
308
309
310
311
312
313
314
315
316
317
318
319
320
321
322
323
324
325
326
327
328
329
330
331
332
333
334
335
336
337
338
339
340
341
342
343
344
345
346
347
348
349
350
351
352
353
354
355
356
357
358
359
360
361
362
363
364
365
366
367
368
369
370
371
372
373
374
375
376
377
378
379
380
381
382
383
384
385
386
387
388
389
390
391
392
393
394
395
396
397
398
399
400
401
402
403
404
405
406
407
408
409
410
411
412
413
414
415
416
417
418
419
420
421
422
423
424
425
426
427
428
429
430
431
432
433
434
435
436
437
438
439
440
441
442
443
444
445
446
447
448
449
450
451
452
453
454
455
456
457
458
459
460
/*
	John Walker's Floating Point Benchmark, derived from...

	      Marinchip Interactive Lens Design System
	             John Walker   December 1980

This  program may be used, distributed, and modified freely as
long as the origin information is preserved.

This is a complete optical design raytracing algorithm,
stripped of its user interface and recast into Go.
It not only determines execution speed on an extremely
floating point (including trig function) intensive
real-world application, it checks accuracy on an algorithm
that is exquisitely sensitive to errors. The performance of
this program is typically far more sensitive to changes in
the efficiency of the trigonometric library routines than the
average floating point program.

Implemented in July 2013 by John Nagle (http://www.animats.com).
*/
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	. "math"
)

// Put all exported functions in local namespace

func cot(x float64) float64 {
	return (1.0 / Tan(x))
}

const max_surfaces = 10

/*  Local variables  */

var current_surfaces int16
var paraxial int16

var clear_aperture float64

var aberr_lspher float64
var aberr_osc float64
var aberr_lchrom float64

var max_lspher float64
var max_osc float64
var max_lchrom float64

var radius_of_curvature float64
var object_distance float64
var ray_height float64
var axis_slope_angle float64
var from_index float64
var to_index float64

var spectral_line [9]float64
var s [max_surfaces][5]float64
var od_sa [2][2]float64

var outarr [8]string /* Computed output of program goes here */

var itercount int /* The iteration counter for the main loop
   in the program is made global so that
   the compiler should not be allowed to
   optimise out the loop over the ray
   tracing code. */

const ITERATIONS = 10000

var niter int = ITERATIONS /* Iteration counter */

/* Reference results.  These happen to
   be derived from a run on Microsoft
   Quick BASIC on the IBM PC/AT. */

var refarr0 = "   Marginal ray          47.09479120920   0.04178472683"
var refarr1 = "   Paraxial ray          47.08372160249   0.04177864821"
var refarr2 = "Longitudinal spherical aberration:        -0.01106960671"
var refarr3 = "    (Maximum permissible):                 0.05306749907"
var refarr4 = "Offense against sine condition (coma):     0.00008954761"
var refarr5 = "    (Maximum permissible):                 0.00250000000"
var refarr6 = "Axial chromatic aberration:                0.00448229032"
var refarr7 = "    (Maximum permissible):                 0.05306749907"

var refarr [8]string = [8]string{refarr0, refarr1, refarr2, refarr3, refarr4, refarr5, refarr6, refarr7}

/* The test case used in this program is the design for a 4 inch
   f/12 achromatic telescope objective used as the example in Wyld's
   classic work on ray tracing by hand, given in Amateur Telescope
   Making, Volume 3 (Volume 2 in the 1996 reprint edition). */

var testcase [4][4]float64 = [4][4]float64{
	[4]float64{27.05, 1.5137, 63.6, 0.52},
	[4]float64{-16.68, 1, 0, 0.138},
	[4]float64{-16.68, 1.6164, 36.7, 0.38},
	[4]float64{-78.1, 1, 0, 0},
}

/*           Calculate passage through surface

If the variable paraxial is 1, the trace through the
surface will be done using the paraxial approximations.
Otherwise, the normal trigonometric trace will be done.

This routine takes the following inputs:

radius_of_curvature      Radius of curvature of surface
                         being crossed.  If 0, surface is plane.

object_distance          Distance of object focus from
                         lens vertex.  If 0, incoming
                         rays are parallel and
                         the following must be specified:

ray_height               Height of ray from axis.  Only
                         relevant if OBJECT.DISTANCE == 0

axis_slope_angle         Angle incoming ray makes with axis
                         at intercept

from_index               Refractive index of medium being left

to_index                 Refractive index of medium being entered.

The outputs are the following variables:

object_distance          Distance from vertex to object focus
                         after refraction.

axis_slope_angle         Angle incoming ray makes with axis
                         at intercept after refraction.

*/

func transit_surface() {
	var iang float64
	var rang float64     /* Refraction angle */
	var iang_sin float64 /* Incidence angle sin */
	var rang_sin float64 /* Refraction angle sin */
	var old_axis_slope_angle float64
	var sagitta float64

	if paraxial > 0 {
		if radius_of_curvature != 0.0 {
			if object_distance == 0.0 {
				axis_slope_angle = 0.0
				iang_sin = ray_height / radius_of_curvature
			} else {
				iang_sin = ((object_distance -
					radius_of_curvature) / radius_of_curvature) *
					axis_slope_angle
			}

			rang_sin = (from_index / to_index) *
				iang_sin
			old_axis_slope_angle = axis_slope_angle
			axis_slope_angle = axis_slope_angle +
				iang_sin - rang_sin
			if object_distance != 0.0 {
				ray_height = object_distance * old_axis_slope_angle
			}
			object_distance = ray_height / axis_slope_angle
			return
		}
		object_distance = object_distance * (to_index / from_index)
		axis_slope_angle = axis_slope_angle * (from_index / to_index)
		return
	}

	if radius_of_curvature != 0.0 {
		if object_distance == 0.0 {
			axis_slope_angle = 0.0
			iang_sin = ray_height / radius_of_curvature
		} else {
			iang_sin = ((object_distance -
				radius_of_curvature) / radius_of_curvature) *
				Sin(axis_slope_angle)
		}
		iang = Asin(iang_sin)
		rang_sin = (from_index / to_index) *
			iang_sin
		old_axis_slope_angle = axis_slope_angle
		axis_slope_angle = axis_slope_angle +
			iang - Asin(rang_sin)
		sagitta = Sin((old_axis_slope_angle + iang) / 2.0)
		sagitta = 2.0 * radius_of_curvature * sagitta * sagitta
		object_distance = ((radius_of_curvature * Sin(
			old_axis_slope_angle+iang)) *
			cot(axis_slope_angle)) + sagitta
		return
	}

	rang = -Asin((from_index / to_index) *
		Sin(axis_slope_angle))
	object_distance = object_distance * ((to_index *
		Cos(-rang)) / (from_index *
		Cos(axis_slope_angle)))
	axis_slope_angle = -rang
}

/*  Perform ray trace in specific spectral line  */

func trace_line(line int, ray_h float64) {

	var i int16

	object_distance = 0.0
	ray_height = ray_h
	from_index = 1.0

	for i = 1; i <= current_surfaces; i++ {
		radius_of_curvature = s[i][1]
		to_index = s[i][2]
		if to_index > 1.0 {
			to_index = to_index + ((spectral_line[4]-
				spectral_line[line])/
				(spectral_line[3]-spectral_line[6]))*((s[i][2]-1.0)/
				s[i][3])
		}
		transit_surface()
		from_index = to_index
		if i < current_surfaces {
			object_distance = object_distance - s[i][4]
		}
	}
}

/*  Initialise when called the first time  */

func main() {
	lambda.Start(HandleRequest)

	/* Process the number of iterations argument, if one is supplied. */
	// if len(os.Args) > 1 {
	// 	var err error
	// 	niter, err = strconv.Atoi(os.Args[1]) // parse argument
	// 	if err != nil || (niter == 0) {
	// 		fmt.Printf("This is John Walker's floating point accuracy and\n")
	// 		fmt.Printf("performance benchmark program.  You call it with\n")
	// 		fmt.Printf("\nfbench <itercount>\n\n")
	// 		fmt.Printf("where <itercount> is the number of iterations\n")
	// 		fmt.Printf("to be executed.  Archival timings should be made\n")
	// 		fmt.Printf("with the iteration count set so that roughly five\n")
	// 		fmt.Printf("minutes of execution is timed.\n")
	// 		os.Exit(0)
	// 	}
	// }
	/* Gimmick: if the iteration is count is negative,
	   suppress user interaction, allowing batch timing
	   runs. */
	// interactive := true
	if niter < 0 {
		niter = -niter
		// interactive = false
	}

	compute()

}

type InputType struct {
	Itercount int `json:"itercount"`
}

func HandleRequest(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	body := InputType{}
	json.Unmarshal([]byte(request.Body), &body)

	fmt.Printf("parsed", body)
	fmt.Println("Doing Float")

	niter = body.Itercount

	return compute()

}

func compute() (events.LambdaFunctionURLResponse, error) {

	var errors int32
	var od_fline float64
	var od_cline float64

	spectral_line[1] = 7621.0   /* A */
	spectral_line[2] = 6869.955 /* B */
	spectral_line[3] = 6562.816 /* C */
	spectral_line[4] = 5895.944 /* D */
	spectral_line[5] = 5269.557 /* E */
	spectral_line[6] = 4861.344 /* F */
	spectral_line[7] = 4340.477 /* G'*/
	spectral_line[8] = 3968.494 /* H */

	/* Load test case into working array */

	clear_aperture = 4.0
	current_surfaces = 4
	var i int16
	for i = 0; i < current_surfaces; i++ {
		for j := 0; j < 4; j++ {
			{
				s[i+1][j+1] = testcase[i][j]
			}
		}
	}

	// if interactive {
	// 	fmt.Printf("Ready to begin John Walker's floating point accuracy\n")
	// 	fmt.Printf("and performance benchmark.  %d iterations will be made.\n\n",
	// 		niter)

	// 	fmt.Printf("\nMeasured run time in seconds should be divided by %f\n", float64(niter)/1000.0)
	// 	fmt.Printf("to normalise for reporting results.  For archival results,\n")
	// 	fmt.Printf("adjust iteration count so the benchmark runs about five minutes.\n\n")

	// 	fmt.Printf("Press return to begin benchmark:")
	// 	fmt.Scanln() // wait for EOL.
	// }
	starttime := time.Now() // start timing

	/* Perform ray trace the specified number of times. */

	for itercount = 0; itercount < niter; itercount++ {

		for paraxial = 0; paraxial <= 1; paraxial++ {

			/* Do main trace in D light */

			trace_line(4, clear_aperture/2.0)
			od_sa[paraxial][0] = object_distance
			od_sa[paraxial][1] = axis_slope_angle
		}
		paraxial = 0

		/* Trace marginal ray in C */

		trace_line(3, clear_aperture/2.0)
		od_cline = object_distance

		/* Trace marginal ray in F */

		trace_line(6, clear_aperture/2.0)
		od_fline = object_distance

		// Compute aberrations of the design

		/* The longitudinal spherical aberration is just the
		   difference between where the D line comes to focus
		   for paraxial and marginal rays. */
		aberr_lspher = od_sa[1][0] - od_sa[0][0]

		/* The offense against the sine condition is a measure
		   of the degree of coma in the design.  We compute it
		   as the lateral distance in the focal plane between
		   where a paraxial ray and marginal ray in the D line
		   come to focus. */
		aberr_osc = 1.0 - (od_sa[1][0]*od_sa[1][1])/
			(Sin(od_sa[0][1])*od_sa[0][0])

			/* The axial chromatic aberration is the distance between
			   where marginal rays in the C and F lines come to focus. */
		aberr_lchrom = od_fline - od_cline

		// Compute maximum acceptable values for each aberration

		max_lspher = Sin(od_sa[0][1])

		/* Maximum longitudinal spherical aberration, which is
		   also the maximum for axial chromatic aberration.  This
		   is computed for the D line. */
		max_lspher = 0.0000926 / (max_lspher * max_lspher)
		max_lchrom = max_lspher
		max_osc = 0.0025 // Max sine condition offence is constant
	}

	elapsedtime := time.Since(starttime) // timing

	// if interactive {
	// 	fmt.Printf("Stop the timer:\007")
	// 	fmt.Scanln() // wait for EOL.
	// }

	/* Now evaluate the accuracy of the results from the last ray trace */

	outarr[0] = fmt.Sprintf("%15s   %21.11f  %14.11f",
		"Marginal ray", od_sa[0][0], od_sa[0][1])
	outarr[1] = fmt.Sprintf("%15s   %21.11f  %14.11f",
		"Paraxial ray", od_sa[1][0], od_sa[1][1])
	outarr[2] = fmt.Sprintf(
		"Longitudinal spherical aberration:      %16.11f",
		aberr_lspher)
	outarr[3] = fmt.Sprintf(
		"    (Maximum permissible):              %16.11f",
		max_lspher)
	outarr[4] = fmt.Sprintf(
		"Offense against sine condition (coma):  %16.11f",
		aberr_osc)
	outarr[5] = fmt.Sprintf(
		"    (Maximum permissible):              %16.11f",
		max_osc)
	outarr[6] = fmt.Sprintf(
		"Axial chromatic aberration:             %16.11f",
		aberr_lchrom)
	outarr[7] = fmt.Sprintf(
		"    (Maximum permissible):              %16.11f",
		max_lchrom)

	/* Now compare the edited results with the master values from
	   reference executions of this program. */

	errors = 0
	for i = 0; i < 8; i++ {
		if outarr[i] != refarr[i] {
			fmt.Printf("\nError in results on line %d...\n", i+1)
			fmt.Printf("Expected:  \"%s\"\n", refarr[i])
			fmt.Printf("Received:  \"%s\"\n", outarr[i])
			fmt.Printf("(Errors)    ")
			k := len(refarr[i])
			for j := 0; j < k; j++ {
				if refarr[i][j] == outarr[i][j] {
					fmt.Printf(" ")
				} else {
					fmt.Printf("^") // indicate character where data did not compare.
				}
				if refarr[i][j] != outarr[i][j] {
					errors++
				}
			}
			fmt.Printf("\n")
		}
	}

	message := ""
	if errors > 0 {
		plural := ""
		if errors > 1 {
			plural = "s"
		}
		message = fmt.Sprintf("\n%d error%s in results.  This is VERY SERIOUS.\n", errors, plural)
		fmt.Printf(message)
	} else {
		fmt.Printf("\nNo errors in results.\n")
		elapsedsecs := elapsedtime.Seconds()
		message = fmt.Sprintf("Elapsed time for %d iterations: %f seconds.   Time per 1000 iterations: %f seconds.\n",
			niter, elapsedsecs, elapsedsecs/float64(niter)*1000)
		fmt.Printf(message)
	}

	ApiResponse := events.LambdaFunctionURLResponse{Body: message, StatusCode: 200}
	return ApiResponse, nil
}
