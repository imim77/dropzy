import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

import FileDrop from './components/FileDrop'

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
      <FileDrop/>
    </>
  )
}

export default App
