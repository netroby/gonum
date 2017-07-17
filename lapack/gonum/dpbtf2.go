// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

// Dpbtf2 computes the Cholesky factorization of a symmetric positive banded
// matrix ab. The matrix ab is n×n with kd diagonal bands. The Choleksy
// factorization computed is
//  A = U^T * U if ul == blas.Upper
//  A = L * L^T if ul == blas.Lower
// ul also specifies the storage of ab. If ul == blas.Upper, then
// ab is stored as an upper-triangular banded matrix with kd super-diagonals,
// and if ul == blas.Lower, ab is stored as a lower-triangular banded matrix
// with kd sub-diagonals. On exit, the banded matrix U or L is stored in-place
// into ab depending on the value of ul. Dpbtf2 returns whether the factorization
// was successfully completed.
//
// Dpbtf2 is the unblocked version of the algorithm, see Dpbtrf for the blocked
// version. Dpbtf2 is an internal routine, exported for testing purposes.
func (Implementation) Dpbtf2(ul blas.Uplo, n, kd int, ab []float64, ldab int) (ok bool) {
	if ul != blas.Upper && ul != blas.Lower {
		panic(badUplo)
	}
	checkSymBanded(ab, n, kd, ldab)
	if n == 0 {
		return
	}
	bi := blas64.Implementation()
	kld := max(1, ldab-1)
	if ul == blas.Upper {
		for j := 0; j < n; j++ {
			// Compute U(J,J) and test for non positive-definiteness
			ajj := ab[j*ldab]
			if ajj <= 0 {
				return false
			}
			ajj = math.Sqrt(ajj)
			ab[j*ldab] = ajj
			// Compute elements j+1:j+kn of row J and update the trailing submatrix
			// within the band.
			kn := min(kd, n-j-1)
			if kn > 0 {
				bi.Dscal(kn, 1/ajj, ab[j*ldab+1:], 1)
				bi.Dsyr(blas.Upper, kn, -1, ab[j*ldab+1:], 1, ab[(j+1)*ldab:], kld)
			}
		}
		return true
	}
	for j := 0; j < n; j++ {
		// Compute L(J,J) and test for non positive-definiteness
		ajj := ab[j*ldab+kd]
		if ajj <= 0 {
			return false
		}
		ajj = math.Sqrt(ajj)
		ab[j*ldab+kd] = ajj

		// Compute elements J+1:J+KN of column J and update the trailing submatrix
		// within the band.
		kn := min(kd, n-j-1)
		if kn > 0 {
			bi.Dscal(kn, 1/ajj, ab[(j+1)*ldab+kd-1:], kld)
			bi.Dsyr(blas.Lower, kn, -1, ab[(j+1)*ldab+kd-1:], kld, ab[(j+1)*ldab+kd:], kld)
		}
	}
	return true
}