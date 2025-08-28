import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'node:path'
import autoprefixer from 'autoprefixer'

const SRC_DIR = path.resolve(__dirname, 'src')
const ASSETS_SRC_DIR = path.resolve(SRC_DIR, 'assets')

// Normalize to POSIX (forward slashes) for Rollup's emitted paths
const toPosix = (p) => p.split(path.sep).join('/')

export default defineConfig(() => {
  return {
    plugins: [vue()],
    base: './',
    css: {
      postcss: {
        plugins: [
          autoprefixer({}), // add options if needed
        ],
      },
    },
    resolve: {
      alias: [
        { find: /^~(.*)$/, replacement: '$1' },
        { find: '@/',
          replacement: `${path.resolve(__dirname, 'src')}/`,
        },
        { find: '@', replacement: path.resolve(__dirname, '/src') },
      ],
      extensions: [
        '.mjs', '.js', '.ts', '.jsx', '.tsx',
        '.json', '.vue', '.scss'
      ],
    },
    build: {
      outDir: "../backend/static/admin",
      emptyOutDir: true,
      target: "es2015",
      sourcemap: true,
      commonjsOptions: {
        include: [/node_modules/],
      },
      rollupOptions: {
        output: {
          // Add hashes for cache busting
          entryFileNames: 'assets/[name].[hash].js',
          chunkFileNames: 'assets/[name].[hash].js',

          assetFileNames: (assetInfo) => {
            const abs = assetInfo?.name || ''
            const ext = path.extname(abs) || ''
            const base = path.basename(abs, ext)

            // If the asset comes from src/assets/**, mirror its subfolder structure
            let relDir = ''
            if (abs) {
              const fromAssets = path.relative(ASSETS_SRC_DIR, path.dirname(abs))
              if (!fromAssets.startsWith('..') && !path.isAbsolute(fromAssets)) {
                relDir = toPosix(fromAssets) // e.g. img/icons
              }
            }

            // If not from src/assets, fall back to grouping by type
            if (!relDir) {
              const extNoDot = ext.replace('.', '').toLowerCase()
              const bucket =
                ['png','jpg','jpeg','gif','webp','svg','avif'].includes(extNoDot) ? 'img' :
                ['woff','woff2','ttf','otf','eot'].includes(extNoDot) ? 'fonts' :
                ['mp4','webm','ogg','mp3','wav','flac'].includes(extNoDot) ? 'media' :
                extNoDot === 'css' ? 'css' :
                ''
              relDir = bucket ? bucket : ''
            }

            // Append hash to filenames
            return relDir
              ? `assets/${relDir}/${base}.[hash][extname]`
              : `assets/${base}.[hash][extname]`
          },

          manualChunks: {
            // define custom chunk splitting here if needed
          },
        },
      },
    },
    server: {
      port: 3000,
      proxy: getProxyOptions()
    },
  }
})

function getProxyOptions() {
  const webserver_port = 8080

  return {
    "^/(app|login|api|assets|files|private)": {
      target: `http://127.0.0.1:${webserver_port}`,
      ws: true,
      rewrite: (path) => path.replace(/^\//, ''),
    },
  }
}
