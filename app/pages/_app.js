import "../styles/globals.css";
function MyApp({ Component, pageProps }) {
  return <Component {...pageProps} />;
}

export default MyApp;

// export default function MyApp({
//   Component,
//   pageProps: { session, ...pageProps },
// }) {
//   return (
//     <SessionProvider session={session}>
//       {" "}
//       <Component {...pageProps} />{" "}
//     </SessionProvider>
//   );
// }
