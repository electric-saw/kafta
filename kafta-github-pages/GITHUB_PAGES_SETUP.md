# GitHub Pages Setup for Kafta

## 📋 Setup Instructions

### 1. Copy Files to Repository
Extract all files from this ZIP to your Kafta repository root directory.

### 2. Repository Structure
Your repository should look like this:
```
kafta/
├── _config.yml                 # Jekyll configuration
├── index.html                  # Homepage
├── _layouts/
│   └── default.html            # Base layout
├── _includes/
│   ├── header.html             # Site header
│   └── footer.html             # Site footer
├── assets/
│   ├── css/
│   │   └── style.css           # Main stylesheet
│   └── js/
│       └── main.js             # JavaScript functionality
├── pages/
│   └── installation.html       # Installation guide
├── img/
│   └── kafta.png               # Logo (use existing)
└── ... (other project files)
```

### 3. Enable GitHub Pages
1. Go to your repository on GitHub
2. Click **Settings** tab
3. Scroll down to **Pages** section
4. Under **Source**, select **Deploy from a branch**
5. Choose **main** branch and **/ (root)** folder
6. Click **Save**

### 4. Wait for Deployment
- GitHub will automatically build and deploy your site
- This usually takes 5-10 minutes
- You'll receive an email when it's ready

### 5. Access Your Site
Your site will be available at:
`https://electric-saw.github.io/kafta`

## 🎨 Customization

### Colors and Branding
Edit `assets/css/style.css` to change:
- Primary color: `--primary-color`
- Secondary color: `--secondary-color`
- Accent color: `--accent-color`

### Content Updates
- **Homepage**: Edit `index.html`
- **Installation**: Edit `pages/installation.html`
- **Navigation**: Edit `_includes/header.html`
- **Footer**: Edit `_includes/footer.html`

### Add New Pages
1. Create new file in `pages/` directory
2. Add Jekyll front matter:
   ```yaml
   ---
   layout: default
   title: "Your Page Title"
   ---
   ```
3. Update navigation in `_includes/header.html`

## 🔧 Features Included

### Responsive Design
- Mobile-friendly navigation
- Responsive grid layouts
- Touch-friendly interactions

### Interactive Elements
- Tabbed installation guides
- Copy-to-clipboard functionality
- Animated counters
- Smooth scrolling

### SEO Optimized
- Meta tags and descriptions
- Open Graph support
- Schema markup
- Sitemap generation

### Performance
- Optimized CSS and JavaScript
- Fast loading times
- Minimal dependencies

## 🚀 Going Live

### Custom Domain (Optional)
1. Add `CNAME` file to repository root
2. Content: `your-domain.com`
3. Configure DNS with your domain provider

### Analytics (Optional)
Add to `_config.yml`:
```yaml
google_analytics: "UA-XXXXXXXXX-X"
```

### Contact Forms (Optional)
Consider services like:
- Netlify Forms
- Formspree
- Google Forms

## 📞 Support

If you need help with GitHub Pages:
- [GitHub Pages Documentation](https://docs.github.com/en/pages)
- [Jekyll Documentation](https://jekyllrb.com/docs/)
- [GitHub Community Forum](https://github.community/)

---
**Happy coding! 🎉**