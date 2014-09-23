package util

// const BUFFER_SIZE = 1024 * 1024

// func StreamBytes(r io.ReadCloser, fn func([]byte) error) error {
//     defer r.Close()
//     buffer := make([]byte, BUFFER_SIZE)

//     for {
//         n, err := io.ReadFull(r, buffer)

//         if err != nil {
//             return err
//         } else if err = fn(buffer); err != nil {
//             return err
//         } else if n < BUFFER_SIZE {
//             return nil
//         }
//     }
// }
