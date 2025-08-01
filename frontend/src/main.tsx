import React from 'react'
import {createRoot} from 'react-dom/client'
import './style.scss'
import {Client} from './client/Client'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
		<Client/>
    </React.StrictMode>
)
