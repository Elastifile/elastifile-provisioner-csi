package retry_test
//
//import (
//	"fmt"
//	"log"
//	"math/rand"
//	"time"
//
//	. "tester"
//
//	"github.com/go-errors/errors"
//	. "github.com/onsi/gomega"
//
//	"retry"
//)
//
//// We want the test to run fast but still do some real sleeping
//const timeout = 50 * time.Millisecond
//
//// const second = time.Second
//
//const second = 10 * time.Millisecond // simulated fast time
//
//type temporaryError struct {
//	message string
//}
//
//func (e *temporaryError) Error() string {
//	return e.message
//}
//
//func (e *temporaryError) Temporary() bool {
//	return true
//}
//
//var _ = Describe("Do", func() {
//
//	It("should succeed when running a successful action", func() {
//		err := retry.Do(timeout,
//			func() error {
//				log.Print("action")
//				return nil
//			})
//		Expect(err).NotTo(HaveOccurred())
//	})
//
//	It("should succeed after retrying a temporarily failed action", func() {
//		retryCount := 0
//		err := retry.Do(timeout,
//			func() error {
//				retryCount++
//				log.Printf("action: try %v", retryCount)
//				if retryCount < 5 {
//					err := &temporaryError{fmt.Sprintf("failed action at try %v", retryCount)}
//					log.Print(err)
//					return err
//				}
//				return nil
//			})
//		Expect(err).NotTo(HaveOccurred())
//	})
//
//	It("should fail after retrying a permanently failed action", func() {
//
//		retryCount := 0
//		err := retry.Do(timeout,
//			func() error {
//				retryCount++
//				log.Printf("action: try %v", retryCount)
//				if retryCount < 5 {
//					err := errors.New("permanent error")
//					log.Print(err)
//					return err
//				}
//				return nil
//			})
//		Expect(err).To(HaveOccurred())
//	})
//
//	It("should fail after timeout", func() {
//		err := retry.Do(timeout,
//			func() error {
//				log.Print("action")
//				time.Sleep(timeout * 2)
//				return nil
//			})
//		Expect(retry.IsTimeout(err)).To(
//			BeTrue(),
//			fmt.Sprintf("Expected to fail after: %s", timeout),
//		)
//	})
//})
//
//var _ = Describe("BasicRetrier", func() {
//	It("Should run BasicRetrier.Do() 5 times 5 seconds each", func() {
//		i := 0
//		start := time.Now()
//		err := retry.Basic{
//			Timeout: 5 * second,
//			Retries: 5,
//		}.Do(func() error {
//			log.Printf("Time elapsed %s\n", time.Since(start))
//			// Imitate work done by action
//			n := rand.Intn(5)
//			time.Sleep(time.Duration(n) * second)
//			i++
//			if i < 5 {
//				return &temporaryError{"Still trying"}
//			}
//			return nil
//		})
//		Expect(err).NotTo(HaveOccurred())
//		Expect(time.Since(start)).To(
//			BeNumerically(">=", 20*second),
//			"BasicRetrier finished too early",
//		)
//		Expect(time.Since(start)).To(
//			BeNumerically("<=", 26*second),
//			"BasicRetrier finished too late",
//		)
//	})
//})
//
//var _ = Describe("LinearRetrier", func() {
//	It("Should run LinearRetrier.Do() 5 times increasing timeout by 2 seconds", func() {
//		i := 0
//		start := time.Now()
//		err := retry.Linear{
//			Timeout: 2 * second,
//			Retries: 5,
//		}.Do(func() error {
//			log.Printf("Time elapsed %s\n", time.Since(start))
//			// Imitate work done by action
//			n := rand.Intn(2)
//			time.Sleep(time.Duration(n) * second)
//			i++
//			if i < 5 {
//				return &temporaryError{"Still trying"}
//			}
//			return nil
//		})
//		Expect(err).NotTo(HaveOccurred())
//		// 2 * sum(1..4) = 2 * 4 * (4 + 1) / 2 = 2 * 10 = 20
//		Expect(time.Since(start)).To(
//			BeNumerically(">=", 20*second),
//			"LinearRetrier finished too early",
//		)
//		Expect(time.Since(start)).To(
//			BeNumerically("<=", 23*second),
//			"LinearRetrier finished too late",
//		)
//	})
//})
//
//var _ = Describe("SigmoidRetrier", func() {
//	It("Should run SigmoidRetrier.Do() 10 times with maximum timeout 5 seconds", func() {
//		i := 0
//		start := time.Now()
//		err := retry.Sigmoid{
//			Limit:   5 * second,
//			Retries: 10,
//		}.Do(func() error {
//			log.Printf("Time elapsed %s\n", time.Since(start))
//			// Imitate work done by action
//			n := rand.Intn(2)
//			time.Sleep(time.Duration(n) * second)
//			i++
//			if i < 10 {
//				return &temporaryError{"Still trying"}
//			}
//			return nil
//		})
//		Expect(err).NotTo(HaveOccurred())
//		//           1
//		// S(t) = -------
//		//             -t
//		//        1 - e
//		// Typical increments:
//		// 0, 0.1, 0.25, 0.6, 1.3, 2.5, 3.6, 4.4, 4.75, 4.9
//		Expect(time.Since(start)).To(
//			BeNumerically(">=", 21*second),
//			"SigmoidRetrier finished too early",
//		)
//		Expect(time.Since(start)).To(
//			BeNumerically("<=", 26*second),
//			"SigmoidRetrier finished too late",
//		)
//	})
//})
//
//var _ = Describe("RootRetrier", func() {
//	It("Should run RootRetrier.Do() 5 times increasing by 5 seconds", func() {
//		i := 0
//		start := time.Now()
//		err := retry.Root{
//			Increment: 5 * second,
//			Retries:   5,
//		}.Do(func() error {
//			log.Printf("Time elapsed %s\n", time.Since(start))
//			// Imitate work done by action
//			n := rand.Intn(2)
//			time.Sleep(time.Duration(n) * second)
//			i++
//			if i < 5 {
//				return &temporaryError{"Still trying"}
//			}
//			return nil
//		})
//		Expect(err).NotTo(HaveOccurred())
//		// 5 * [sqrt(1), sqrt(2), sqrt(3), sqrt(4)] ~ 5 * 6.2 ~ 31
//		// 5 * [sqrt(1), sqrt(2), sqrt(3), sqrt(4), sqrt(5)] ~ 5 * 8.4 ~ 40
//		Expect(time.Since(start)).To(
//			BeNumerically(">=", 30*second),
//			"RootRetrier finished too early",
//		)
//		Expect(time.Since(start)).To(
//			BeNumerically("<=", 41*second),
//			"RootRetrier finished too late",
//		)
//	})
//})
//
//var _ = Describe("OpportunisticRetrier", func() {
//	It("Should run OpportunisticRetrier.Do() 5 times", func() {
//		i := 0
//		start := time.Now()
//		window := time.Now()
//		err := retry.Opportunistic{
//			Total:   30 * second,
//			Retries: 5,
//		}.Do(func() error {
//			log.Printf("Time elapsed %s\n", time.Since(start))
//			// Imitate work done by action
//			n := rand.Intn(2)
//			window = window.Add(time.Duration(n) * second)
//			time.Sleep(time.Duration(n) * second)
//			i++
//			if i < 5 {
//				return &temporaryError{"Still trying"}
//			}
//			return nil
//		})
//		Expect(err).NotTo(HaveOccurred())
//		delta := time.Since(window)
//		Expect(delta).To(
//			BeNumerically("<=", 1*second),
//			"OpportunisticRetrier finished too late",
//		)
//	})
//})
