export default function HomePage() {
  const videoRef = useRef<any>();

  useEffect(() => {
    timerRef.current = setTimeout(() => {
      setFold(true);
    }, 5000);
  }, []);
}
